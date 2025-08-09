# Operator Workaround for PBSS Init Issue

## Quick Fix for Node Operators

If you're experiencing the "Failed to truncate extra state histories" error when restarting your geth node with PBSS, here's how to work around it:

### Recommended: Avoid Repeat Init

**Instead of:**
```bash
# DON'T DO THIS on every restart
geth init --datadir /data genesis.json
geth --datadir /data --syncmode full ...
```

**Do this:**
```bash
# First time only
geth init --datadir /data genesis.json

# All subsequent restarts (no init)
geth --datadir /data --syncmode full ...
```

### For Automation/Scripts

Modify your startup scripts to check if the database already exists:

```bash
#!/bin/bash
DATADIR="/data"
GENESIS="genesis.json"

# Only init if chaindata doesn't exist
if [ ! -d "$DATADIR/geth/chaindata" ]; then
    echo "Initializing new database..."
    geth init --datadir "$DATADIR" "$GENESIS"
else
    echo "Database already exists, skipping init"
fi

# Start geth normally
geth --datadir "$DATADIR" --syncmode full ...
```

### If You Need to Update Chain Config

When you need to update the chain configuration (e.g., for a hard fork):

1. **Stop the node**
2. **Backup your data** (especially `/geth/chaindata`)
3. **Try starting without init first** - geth can often handle config updates automatically
4. **If that fails**, use hash scheme temporarily:
   ```bash
   geth --state.scheme=hash init --datadir /data new-genesis.json
   ```

### Emergency Recovery

If your node is already in a failed state:

```bash
# Option 1: Reset state histories (keeps blockchain data)
rm -rf /data/geth/chaindata/ancient/state
geth init --datadir /data genesis.json

# Option 2: Full reset (last resort)
rm -rf /data/geth/chaindata
geth init --datadir /data genesis.json
```

## Why This Happens

- **Upstream geth v1.16.2 bug** with path-based state scheme
- Running `init` multiple times corrupts state history metadata
- The issue does NOT occur with hash-based state scheme

## Long-term Solution

Apply the patch in `minimal_pbss_init_fix.patch` to your geth build, or wait for an upstream fix.

## Testing Your Setup

After implementing the workaround:

1. Start your node normally
2. Stop it gracefully (`kill -SIGTERM` or Ctrl+C)
3. Restart without init
4. Verify it starts without errors

## Support

If you continue experiencing issues:
1. Check you're not running init on restart
2. Verify your automation scripts follow the pattern above
3. Consider using `state.scheme=hash` until the bug is fixed upstream
