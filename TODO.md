     -------Z
    /      /
---A------C---M
    \ \      /
     \ o----B
      \
---R---W


                merge
                ---
                   moveHead
                   ----
                        wasFF
                       -------
CC
    (default) -> C, no, no-ff  // Already up to date.
    --ff-only -> C, no, no-ff  // Already up to date.
    --no-ff   -> C, no, no-ff  // Already up to date.
CA
    (default) -> C, no, no-ff  // Already up to date.
    --ff-only -> C, no, no-ff  // Already up to date.
    --no-ff   -> C, no, no-ff  // Already up to date.
AC
    (default) -> C, yes, ff    // Updating HASH-1..HASH-2 \n Fast-forward (no commit created; -m option ignored)
    --ff-only -> C, yes, ff    // Updating HASH-1..HASH-2 \n Fast-forward (no commit created; -m option ignored)
    --no-ff   -> Z, yes, no-ff // Merge made by the 'recursive' strategy.
BC
    (default) -> M, yes, no-ff // Merge made by the 'recursive' strategy.
    --ff-only -> ERROR         // 128 Not possible to fast-forward, aborting.
    --no-ff   -> M, yes, no-ff // Merge made by the 'recursive' strategy.
AR
    (default) -> ERROR         // 128 Refusing to merge unrelated histories
    --ff-only -> ERROR         // 128 Refusing to merge unrelated histories
    --no-ff   -> ERROR         // 128 Refusing to merge unrelated histories
    --allow-unrelated-histories
              -> W, yes, no-ff // Merge made by the 'recursive' strategy.

switch merge
    case 1st: Already up to date.
    case 2nd: Updating HASH-1..HASH-2 \n Fast-forward
    default:  Merge made by the 'recursive' strategy.

go-git merge target
    -> Merge HEAD target
        If HEAD is passed, it is updated.
        To do so:
        -> merge HEAD target
            produces the merge, and done
            needs a clear WorkingDirectory

