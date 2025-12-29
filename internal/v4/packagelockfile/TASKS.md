# Package Lock File Format Migration Tasks

## Overview
Migrate from nested registry/package structure to flat "registry/package" key format.

## Tasks

### Core Implementation
- [x] 1. Update PackageLockfile struct - Change Packages field type
- [x] 2. Update GetPackageLockInfo method - Use flat key lookup
- [x] 3. Update UpsertPackageLockInfo method - Remove nested map logic
- [x] 4. Update UpdatePackageName method - Handle flat keys
- [x] 5. Update RemovePackageLockInfo method - Simplify deletion
- [x] 6. Update UpdateRegistryName method - Iterate and rename keys
- [x] 7. Update readLockFile method - Remove nested map initialization

### Testing
- [x] 8. Update TestFileManager_GetPackageLockInfo - Change test data structure
- [x] 9. Update TestFileManager_GetPackageLockfile - Change expected structure
- [x] 10. Update TestFileManager_UpsertPackageLockInfo - Change test data and assertions
- [x] 11. Update TestFileManager_UpdatePackageName - Change test data and logic
- [x] 12. Update TestFileManager_RemovePackageLockInfo - Change test data and assertions
- [x] 13. Update TestFileManager_UpdateRegistryName - Change test data and logic
- [x] 14. Update createTestLockfile helper - Change structure
- [x] 15. Verify all test cases pass

### Validation
- [x] 16. Test with sample lock file format
- [x] 17. Ensure backward compatibility handling (if needed)
- [x] 18. Add helper functions to core package
- [x] 19. Update packagelockfile to use core helpers
- [x] 20. Update manifest to use core helpers

## Progress
- [x] Task file created
- [x] Implementation completed (tasks 1-7)
- [x] Tests updated (tasks 8-15)
- [x] All tests passing
- [x] Validation completed (tasks 16-20)
- [x] Helper functions centralized in core
- [x] Migration complete