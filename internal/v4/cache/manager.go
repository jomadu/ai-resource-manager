package cache

// DIRECTORY STRUCTURE:
// ~/.arm/
// ├── .armrc                                # ARM configuration (optional)
// └── registries/
//     └── {registry-key}/
//         ├── metadata.json                 # Registry metadata + timestamps
//         ├── repo/                         # Git clone (git registries only)
//         │   ├── .git/
//         │   └── source-files...
//         └── packages/
//             └── {package-key}/
//                 ├── metadata.json         # Package metadata + timestamps
//                 └── {version}/
//                     ├── metadata.json     # Version metadata + timestamps
//                     └── files/
//                         └── extracted-files...
//
// Cache provides two main components:
// - RegistryCache: handles registry directory and metadata.json
// - PackageCache: handles package storage within registry
