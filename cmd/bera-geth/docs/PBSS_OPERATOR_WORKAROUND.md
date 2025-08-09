# Operator Workaround for PBSS Init Issue

## Quick Fix for Node Operators

If you're experiencing the "Failed to truncate extra state histories" error when restarting your geth node with PBSS, here's how to work around it:

### Recommended: Avoid Repeat Init

**Instead of:**
```bash
# DON'T DO THIS on every restart
bera-geth init --datadir /data eth-genesis.json
bera-geth --datadir /data --syncmode full ...
```

**Do this:**
```bash
# First time only
bera-geth init --datadir /data eth-genesis.json

# All subsequent restarts (no init)
bera-geth --datadir /data --syncmode full ...
```

### For Automation/Scripts

Modify your startup scripts to check if the database already exists:

```bash
#!/bin/bash
DATADIR="/data"
GENESIS="eth-genesis.json"

# Only init if chaindata doesn't exist
if [ ! -d "$DATADIR/bera-geth/chaindata" ]; then
    echo "Initializing new database..."
    bera-geth init --datadir "$DATADIR" "$GENESIS"
else
    echo "Database already exists, skipping init"
fi

# Start geth normally
bera-geth --datadir "$DATADIR" --syncmode full ...
```

### If You Need to Update Chain Config

When you need to update the chain configuration (e.g., for a hard fork):

1. **Stop the node**
2. **Backup your data** (especially `/bera-geth/chaindata`)
3. **Try starting without init first** - geth can often handle config updates automatically
4. **If that fails**, use hash scheme temporarily:
   ```bash
   bera-geth --state.scheme=hash init --datadir /data new-genesis.json
   ```

### Emergency Recovery

If your node is already in a failed state:

```bash
# Option 1: Reset state histories (keeps blockchain data)
rm -rf /data/bera-geth/chaindata/ancient/state
bera-geth init --datadir /data eth-genesis.json

# Option 2: Full reset (last resort)
rm -rf /data/bera-geth/chaindata
bera-geth init --datadir /data eth-genesis.json
```

## Why This Happens

- **Upstream geth v1.16.2 bug** with path-based state scheme
- Running `init` multiple times corrupts state history metadata
- The issue does NOT occur with hash-based state scheme


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
