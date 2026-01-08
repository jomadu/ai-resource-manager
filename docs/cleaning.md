# Cleaning Utilities

## Clean Cache

Clean the local cache directory:
```bash
arm clean cache
```

Aggressive cleanup (remove all cached data):
```bash
arm clean cache --nuke
```

## Clean Sinks

Clean sink directories based on ARM index:
```bash
arm clean sinks
```

Complete cleanup (remove entire ARM directory):
```bash
arm clean sinks --nuke
```