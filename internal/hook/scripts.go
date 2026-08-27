package hook

import _ "embed"

const commitMsgHookName = "commit-msg"

// defaultCommitMsgHookScript is the shell script installed into .git/hooks/commit-msg.
// It invokes commitkit check on the message file passed by Git.
//
//go:embed scripts/commit-msg.sh
var defaultCommitMsgHookScript []byte
