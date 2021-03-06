package git

import (
	"errors"

	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp/sideband"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

// SubmoduleRescursivity defines how depth will affect any submodule recursive
// operation.
type SubmoduleRescursivity uint

const (
	// DefaultRemoteName name of the default Remote, just like git command.
	DefaultRemoteName = "origin"

	// NoRecurseSubmodules disables the recursion for a submodule operation.
	NoRecurseSubmodules SubmoduleRescursivity = 0
	// DefaultSubmoduleRecursionDepth allow recursion in a submodule operation.
	DefaultSubmoduleRecursionDepth SubmoduleRescursivity = 10
)

var (
	ErrMissingURL = errors.New("URL field is required")
)

// CloneOptions describes how a clone should be performed.
type CloneOptions struct {
	// The (possibly remote) repository URL to clone from.
	URL string
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Name of the remote to be added, by default `origin`.
	RemoteName string
	// Remote branch to clone.
	ReferenceName plumbing.ReferenceName
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Limit fetching to the specified number of commits.
	Depth int
	// RecurseSubmodules after the clone is created, initialize all submodules
	// within, using their default settings. This option is ignored if the
	// cloned repository does not have a worktree.
	RecurseSubmodules SubmoduleRescursivity
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
}

// Validate validates the fields and sets the default values.
func (o *CloneOptions) Validate() error {
	if o.URL == "" {
		return ErrMissingURL
	}

	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if o.ReferenceName == "" {
		o.ReferenceName = plumbing.HEAD
	}

	return nil
}

// PullOptions describes how a pull should be performed.
type PullOptions struct {
	// Name of the remote to be pulled. If empty, uses the default.
	RemoteName string
	// Remote branch to clone.  If empty, uses HEAD.
	ReferenceName plumbing.ReferenceName
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Limit fetching to the specified number of commits.
	Depth int
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// RecurseSubmodules controls if new commits of all populated submodules
	// should be fetched too.
	RecurseSubmodules SubmoduleRescursivity
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
}

// Validate validates the fields and sets the default values.
func (o *PullOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if o.ReferenceName == "" {
		o.ReferenceName = plumbing.HEAD
	}

	return nil
}

// FetchOptions describes how a fetch should be performed
type FetchOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	RefSpecs   []config.RefSpec
	// Depth limit fetching to the specified number of commits from the tip of
	// each remote branch history.
	Depth int
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
}

// Validate validates the fields and sets the default values.
func (o *FetchOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	for _, r := range o.RefSpecs {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// PushOptions describes how a push should be performed.
type PushOptions struct {
	// RemoteName is the name of the remote to be pushed to.
	RemoteName string
	// RefSpecs specify what destination ref to update with what source
	// object. A refspec with empty src can be used to delete a reference.
	RefSpecs []config.RefSpec
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
}

// Validate validates the fields and sets the default values.
func (o *PushOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if len(o.RefSpecs) == 0 {
		o.RefSpecs = []config.RefSpec{
			config.RefSpec(config.DefaultPushRefSpec),
		}
	}

	for _, r := range o.RefSpecs {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// SubmoduleUpdateOptions describes how a submodule update should be performed.
type SubmoduleUpdateOptions struct {
	// Init, if true initializes the submodules recorded in the index.
	Init bool
	// NoFetch tell to the update command to not fetch new objects from the
	// remote site.
	NoFetch bool
	// RecurseSubmodules the update is performed not only in the submodules of
	// the current repository but also in any nested submodules inside those
	// submodules (and so on). Until the SubmoduleRescursivity is reached.
	RecurseSubmodules SubmoduleRescursivity
}

// CheckoutOptions describes how a checkout 31operation should be performed.
type CheckoutOptions struct {
	// Hash to be checked out, if used HEAD will in detached mode. Branch and
	// Hash are mutual exclusive.
	Hash plumbing.Hash
	// Branch to be checked out, if Branch and Hash are empty is set to `master`.
	Branch plumbing.ReferenceName
	// Force, if true when switching branches, proceed even if the index or the
	// working tree differs from HEAD. This is used to throw away local changes
	Force bool
}

// Validate validates the fields and sets the default values.
func (o *CheckoutOptions) Validate() error {
	if o.Branch == "" {
		o.Branch = plumbing.Master
	}

	return nil
}

// ResetMode defines the mode of a reset operation.
type ResetMode int8

const (
	// HardReset resets the index and working tree. Any changes to tracked files
	// in the working tree are discarded.
	HardReset ResetMode = iota
	// MixedReset resets the index but not the working tree (i.e., the changed
	// files are preserved but not marked for commit) and reports what has not
	// been updated. This is the default action.
	MixedReset
	// MergeReset resets the index and updates the files in the working tree
	// that are different between Commit and HEAD, but keeps those which are
	// different between the index and working tree (i.e. which have changes
	// which have not been added).
	//
	// If a file that is different between Commit and the index has unstaged
	// changes, reset is aborted.
	MergeReset
)

// ResetOptions describes how a reset operation should be performed.
type ResetOptions struct {
	// Commit, if commit is pressent set the current branch head (HEAD) to it.
	Commit plumbing.Hash
	// Mode, form resets the current branch head to Commit and possibly updates
	// the index (resetting it to the tree of Commit) and the working tree
	// depending on Mode. If empty MixedReset is used.
	Mode ResetMode
}

// Validate validates the fields and sets the default values.
func (o *ResetOptions) Validate(r *Repository) error {
	if o.Commit == plumbing.ZeroHash {
		ref, err := r.Head()
		if err != nil {
			return err
		}

		o.Commit = ref.Hash()
	}

	return nil
}

// LogOptions describes how a log action should be performed.
type LogOptions struct {
	// When the From option is set the log will only contain commits
	// reachable from it. If this option is not set, HEAD will be used as
	// the default From.
	From plumbing.Hash
}
