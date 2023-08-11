# Object Graph Manifest for Velero

## Abstract

One to two sentences that describes the goal of this proposal and the problem being solved by the proposed change.
The reader should be able to tell by the title, and the opening paragraph, if this document is relevant to them.

Currently, Velero does not have a complete manifest of everything in the backup, aside from the backup tarball itself.
This change introduces a new data structure to be stored with a backup in object storage which will allow for more efficient operations in reporting of what a backup contains.
Additionally, this manifest should enable advancements in Velero's features and architecture, enabling dry-run support, concurrent backup and restore operations, and reliable restoration of complex applications.

## Background

Right now, Velero backs up items one at a time, sorted by API Group and namespace.
It also restores items one at a time, using the restoreResourcePriorities flag to indicate which order API Groups should have their objects restored first.
While this does work currently, it presents challenges for more complex applications that have their dependencies in the form of a graph rather than strictly linear.

For example, Cluster API clusters are a set of complex Kubernetes objects that require that the "root" objects are restored first, before their "leaf" objects.
If a Cluster that a ClusterResourceSetBinding refers to does not exist, then a restore of the CAPI cluster will fail.

Additionally, Velero does not have a reliable way to communicate what objects will be affected in a backup or restore operation without actually performing the operation.
This complicates dry-run tasks, because a user must simply perform the action without knowing what will be touched.
It also complicates allowing backups and restores to run in parallel, because there is currently no way to know if a single Kubernetes object is included in multiple backups or restores, which can lead to unreliability, deadlocking, and race conditions were Velero made to be more concurrent today.

## Goals

- Introduce a manifest data structure that defines the contents of a backup.
- Store the manifest data into object storage alongside existing backup data.

## Non Goals

This proposal seeks to enable, but not define, the following.

- Implementing concurrency beyond what already exists in Velero.
- Implementing a dry-run feature.
- Implementing a new restore ordering procedure.

While the data structure should take these scenarios into account, they will not be implemented alongside it.

## High-Level Design

To uniquely identify a Kubernetes object within a cluster or backup, the following fields are sufficient:

- API Group and Version (example: backup.velero.io/v1)
- Namespace 
- Name
- Labels

This criteria covers the majority of Velero's inclusion or exclusion logic.
However, some additional fields enable further use cases.

- Owners, which are other Kubernetes objects that have some relationship to this object. They may be strict or soft dependencies.
- Annotations, which provide extra metadata about the object that might be useful for other programs to consume.
- UUID generated by Kubernetes. This is useful in defining Owner relationships, providing a single, immutable key to find an object. This is _not_ considered at restore time, only internally for defining links.

All of this information already exists within a Velero backup's tarball of resources, but extracting such data is inefficient.
The entire tarball must be downloaded and extracted, and then JSON within parsed to read labels, owners, annotations, and a UUID.
The rest of the information is encoded in the file system structure within the Velero backup tarball.
While doable, this is heavyweight in terms of time and potentially memory.

Instead, this proposal suggests adding a new manifest structure that is kept alongside the backup tarball.
This structure would contain the above fields only, and could be used to perform inclusion/exclusion logic on a backup, select a resource from within a backup, and do set operations over backup or restore contents to identify overlapping resources.

Here are some use cases that this data structure should enable, that have been difficult to implement prior to its existence:

- A dry-run operation on backup, informing the user what would be selected if they were to perform the operation.
 A manifest could be created and saved, allowing for a user to do a dry-run, then accept it to perform the backup.
 Restore operations can be treated similarly.
- Efficient, non-overlapping parallelization of backup and restore operations.
 By building or reading a manifest before performing a backup or restore, Velero can determine if there are overlapping resources.
 If there are no overlaps, the operations can proceed in parallel.
 If there are overlaps, the operations can proveed serially.
- Graph-based restores for non-linear dependencies.
 Not all resources in a Kubernetes cluster can be defined in a strict, linear way.
 They may have multiple owners, and writing BackupItemActions or RestoreItemActions to simply return a chain of owners is not an efficient way to support the many Kubernetes operators/controllers being written.
 Instead, by having a manifest with enough information, Velero can build a discrete list that ensures dependencies are restored before their dependents, with less input from plugin authors.

## Detailed Design

The Manifest data structure would look like this, in Go type structure:

```golang
// NamespacedItems maps a given namespace to all of its contained items.
type NamespacedItems map[string]*Item

// APIGroupNamespaces maps an API group/version to a map of namespaces and their items.
type KindNamespaces map[string]NamespacedItems

type Manifest struct {
    // Kinds holds the top level map of all resources in a manifest.
    Kinds KindNamespaces

    // Index is used to look up an individual item quickly based on UUID.
    // This enables fetching owners out of the maps more efficiently at the cost of memory space.
    Index map[string]*Item
}


// Item represents a Kubernetes resource within a backup based on it's selectable criteria.
// It is not the whole Kubernetes resource as retrieved from the API server, but rather a collection of important fields needed for filtering.
type Item struct {
    // Kubernetes API group which this Item belongs to.
    // Could be a core resource, or a CustomResourceDefinition.
    APIGroup string

    // Version of the APIGroup that the Item belongs to.
    APIVersion string

    // Kubernetes namespace which contains this item.
    // Empty string for cluster-level resource.
    Namespace string

    // Item's given name.
    Name string

    // Map of labels that the Item had at backup time.
    Labels map[string]string

    // Map of annotations that the Item had at Backup time.
    // Useful for plugins that may decide to process only Items with specific annotations.
    Annotations map[string]string

    // Owners is a list of UUIDs to other items that own or refer to this item.
    Owners []string

    // Manifest is a pointer to the Manifest in which this object is contained.
    // Useful for getting access to things like the Manifest.Index map.
    Manifest *Manifest
}
```

In addition to the new types, the following Go interfaces would be provided for convenience.

```golang
type Itermer interface {
    // Returns the Item as a string, following the current Velero backup version 1.1.0 tarball structure format.
    // <APIGroup>/<Namespace>/<APIVersion>/<name>.json
    String() string

    // Owners returns a slice of realized Items that own or refer to the current Item.
    // Useful for building out a full graph of Items to restore.
    // Will use the UUIDs in Item.Owners to look up the owner Items in the Manifest.
    Owners() []*Item

    // Kind returns the Kind of an object, which is a combination of the APIGroup and APIVersion.
    // Useful for verifying the needed CustomResourceDefinition exists before actually restoring this Item.
    Kind() *Item

    // Children returns a slice of all Items that refer to this item as an Owner.
    Children() []*Items
}

// This error type is being created in order to make reliable sentinel errors.
// See https://dave.cheney.net/2019/06/10/constant-time for more details.
type ManifestError string

func (e ManifestError) Error() string {
    return string(e)
}

const ItemAlreadyExists = ManifestError("item already exists in manifest")

type Manifester interface {
    // Set returns the entire list of resources as a set of strings (using Itemer.String).
    // This is useful for comparing two manifests and determining if they have any overlapping resources.
    // In the future, when implementing concurrent operations, this can be used as a sanity check to ensure resources aren't being backed up or restored by two operations at once.
    Set() sets.String

    // Adds an item to the appropriate APIGroup and Namespace within a Manifest
    // Returns (true, nil) if the Item is successfully added to the Manifest,
    // Returns (false, ItemAlreadyExists) if the Item is already in the Manifest.
    Add(*Item) (bool, error)
}
```

### Serialization

The entire `Manifest` should be serialized into the `manifest.json` file within the object storage for a single backup.
It is possible that this file could also be compressed for space efficiency.

### Memory Concerns

Because the `Manifest` is holding a minimal amount of data, memory sizes should not be a concern for most clusters.
TODO: Document known limits on API group name, resource name, and kind name character limits.

## Security Considerations

Introducing this manifest does not increase the attack surface of Velero, as this data is already present in the existing backups.
Storing the manifest.json file next to the existing backup data in the object storage does not change access patterns.

## Compatibility

The introduction of this file should trigger Velero backup version 1.2.0, but it will not interfere with Velero versions that do not support the `Manifest` as the file will be additive.
In time, this file will replace the `<backupname>-resource-list.json.gz` file, but for compatibility the two will appear side by side.

When first implemented, Velero should simply build the `Manifest` as it backs up items, and serialize it at the end.
Any logic changes that rely on the `Manifest` file must be introduced with their own design document, with their own compatibility concerns.

## Implementation

The `Manifest` object will _not_ be implemented as a Kubernetes CustomResourceDefinition, but rather one of Velero's own internal constructs.

Implementation for the data structure alone should be minimal - the types will need to be defined in a `manifest` package.
Then, the backup process should create a `Manifest`, passing it to the various `*Backuppers` in the `backup` package.
These methods will insert individual `Items` into the `Manifest`.
Finally, logic should be added to the `persistence` package to ensure that the new `manifest.json` file is uploadable and allowed.

## Alternatives Considered

None so far.

## Open Issues

- When should compatibility with the `<backupname>-resource-list.json.gz` file be dropped?
- What are some good test case Kubernetes resources and controllers to try this out with?
Cluster API seems like an obvious choice, but are there others?
- Since it is not implemented as a CustomResourceDefinition, how can a `Manifest` be retained so that users could issue a dry-run command, then perform their actual desire operation?
Could it be stored in Velero's temp directories?
Note that this is making Velero itself more stateful.