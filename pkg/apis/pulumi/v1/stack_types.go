package v1

import (
	"github.com/pulumi/pulumi-kubernetes-operator/pkg/apis/pulumi/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StackStatus defines the observed state of Stack
type StackStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Stack is the Schema for the stacks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=stacks,scope=Namespaced
// +kubebuilder:storageversion
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec          `json:"spec,omitempty"`
	Status shared.StackStatus `json:"status,omitempty"`
}

// StackSpec defines the desired state of Pulumi Stack being managed by this operator.
type StackSpec struct {
	// Auth info:

	// (optional) AccessTokenSecret is the name of a secret containing the PULUMI_ACCESS_TOKEN for Pulumi access.
	// Deprecated: use EnvRefs with a "secret" entry with the key PULUMI_ACCESS_TOKEN instead.
	AccessTokenSecret string `json:"accessTokenSecret,omitempty"`

	// (optional) Envs is an optional array of config maps containing environment variables to set.
	// Deprecated: use EnvRefs instead.
	Envs []string `json:"envs,omitempty"`

	// (optional) EnvRefs is an optional map containing environment variables as keys and stores descriptors to where
	// the variables' values should be loaded from (one of literal, environment variable, file on the
	// filesystem, or Kubernetes secret) as values.
	EnvRefs map[string]shared.ResourceRef `json:"envRefs,omitempty"`

	// (optional) SecretEnvs is an optional array of secret names containing environment variables to set.
	// Deprecated: use EnvRefs instead.
	SecretEnvs []string `json:"envSecrets,omitempty"`

	// (optional) Backend is an optional backend URL to use for all Pulumi operations.<br/>
	// Examples:<br/>
	//   - Pulumi Service:              "https://app.pulumi.com" (default)<br/>
	//   - Self-managed Pulumi Service: "https://pulumi.acmecorp.com" <br/>
	//   - Local:                       "file://./einstein" <br/>
	//   - AWS:                         "s3://<my-pulumi-state-bucket>" <br/>
	//   - Azure:                       "azblob://<my-pulumi-state-bucket>" <br/>
	//   - GCP:                         "gs://<my-pulumi-state-bucket>" <br/>
	// See: https://www.pulumi.com/docs/intro/concepts/state/
	Backend string `json:"backend,omitempty"`

	// Stack identity:

	// Stack is the fully qualified name of the stack to deploy (<org>/<stack>).
	Stack string `json:"stack"`
	// (optional) Config is the configuration for this stack, which can be optionally specified inline. If this
	// is omitted, configuration is assumed to be checked in and taken from the source repository.
	Config map[string]string `json:"config,omitempty"`
	// (optional) Secrets is the secret configuration for this stack, which can be optionally specified inline. If this
	// is omitted, secrets configuration is assumed to be checked in and taken from the source repository.
	// Deprecated: use SecretRefs instead.
	Secrets map[string]string `json:"secrets,omitempty"`

	// (optional) SecretRefs is the secret configuration for this stack which can be specified through ResourceRef.
	// If this is omitted, secrets configuration is assumed to be checked in and taken from the source repository.
	SecretRefs map[string]shared.ResourceRef `json:"secretsRef,omitempty"`
	// (optional) SecretsProvider is used to initialize a Stack with alternative encryption.
	// Examples:
	//   - AWS:   "awskms:///arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34bc-56ef-1234567890ab?region=us-east-1"
	//   - Azure: "azurekeyvault://acmecorpvault.vault.azure.net/keys/mykeyname"
	//   - GCP:   "gcpkms://projects/MYPROJECT/locations/MYLOCATION/keyRings/MYKEYRING/cryptoKeys/MYKEY"
	//   -
	// See: https://www.pulumi.com/docs/intro/concepts/secrets/#initializing-a-stack-with-alternative-encryption
	SecretsProvider string `json:"secretsProvider,omitempty"`

	// Source control: either GitRepo or FluxSource fields should be populated.
	GitRepo   *InlineGitRepo   `json:",inline,omitempty"`
	SourceRef *SourceReference `json:"sourceRef,omitempty"`

	// (optional) RepoDir is the directory to work from in the project's source repository
	// where Pulumi.yaml is located. It is used in case Pulumi.yaml is not
	// in the project source root.
	RepoDir string `json:"repoDir,omitempty"`

	// Lifecycle:

	// (optional) Refresh can be set to true to refresh the stack before it is updated.
	Refresh bool `json:"refresh,omitempty"`
	// (optional) ExpectNoRefreshChanges can be set to true if a stack is not expected to have
	// changes during a refresh before the update is run.
	// This could occur, for example, is a resource's state is changing outside of Pulumi
	// (e.g., metadata, timestamps).
	ExpectNoRefreshChanges bool `json:"expectNoRefreshChanges,omitempty"`
	// (optional) DestroyOnFinalize can be set to true to destroy the stack completely upon deletion of the CRD.
	DestroyOnFinalize bool `json:"destroyOnFinalize,omitempty"`
	// (optional) RetryOnUpdateConflict issues a stack update retry reconciliation loop
	// in the event that the update hits a HTTP 409 conflict due to
	// another update in progress.
	// This is only recommended if you are sure that the stack updates are
	// idempotent, and if you are willing to accept retry loops until
	// all spawned retries succeed. This will also create a more populated,
	// and randomized activity timeline for the stack in the Pulumi Service.
	RetryOnUpdateConflict bool `json:"retryOnUpdateConflict,omitempty"`

	// (optional) UseLocalStackOnly can be set to true to prevent the operator from
	// creating stacks that do not exist in the tracking git repo.
	// The default behavior is to create a stack if it doesn't exist.
	UseLocalStackOnly bool `json:"useLocalStackOnly,omitempty"`

	// (optional) ResyncFrequencySeconds when set to a non-zero value, triggers a resync of the stack at
	// the specified frequency even if no changes to the custom-resource are detected.
	// If branch tracking is enabled (branch is non-empty), commit polling will occur at this frequency.
	// The minimal resync frequency supported is 60 seconds.
	ResyncFrequencySeconds int64 `json:"resyncFrequencySeconds,omitempty"`
}

type InlineGitRepo struct {
	// ProjectRepo is the git source control repository from which we fetch the project code and configuration.
	//+optional
	ProjectRepo string `json:"projectRepo,omitempty"`
	// (optional) GitAuthSecret is the the name of a secret containing an
	// authentication option for the git repository.
	// There are 3 different authentication options:
	//   * Personal access token
	//   * SSH private key (and it's optional password)
	//   * Basic auth username and password
	// Only one authentication mode will be considered if more than one option is specified,
	// with ssh private key/password preferred first, then personal access token, and finally
	// basic auth credentials.
	// Deprecated. Use GitAuth instead.
	GitAuthSecret string `json:"gitAuthSecret,omitempty"`

	// (optional) GitAuth allows configuring git authentication options
	// There are 3 different authentication options:
	//   * SSH private key (and its optional password)
	//   * Personal access token
	//   * Basic auth username and password
	// Only one authentication mode will be considered if more than one option is specified,
	// with ssh private key/password preferred first, then personal access token, and finally
	// basic auth credentials.
	GitAuth *shared.GitAuthConfig `json:"gitAuth,omitempty"`
	// (optional) Commit is the hash of the commit to deploy. If used, HEAD will be in detached mode. This
	// is mutually exclusive with the Branch setting. Either value needs to be specified.
	Commit string `json:"commit,omitempty"`
	// (optional) Branch is the branch name to deploy, either the simple or fully qualified ref name, e.g. refs/heads/master. This
	// is mutually exclusive with the Commit setting. Either value needs to be specified.
	// When specified, the operator will periodically poll to check if the branch has any new commits.
	// The frequency of the polling is configurable through ResyncFrequencySeconds, defaulting to every 60 seconds.
	Branch string `json:"branch,omitempty"`
	// (optional) ContinueResyncOnCommitMatch - when true - informs the operator to continue trying to update stacks
	// even if the commit matches. This might be useful in environments where Pulumi programs have dynamic elements
	// for example, calls to internal APIs where GitOps style commit tracking is not sufficient.
	// Defaults to false, i.e. when a particular commit is successfully run, the operator will not attempt to rerun the
	// program at that commit again.
	ContinueResyncOnCommitMatch bool `json:"continueResyncOnCommitMatch,omitempty"`
}

type SourceReference struct {
	// The API version of the source; e.g., `source.toolkit.fluxcd.io/v1beta2`
	APIVersion string `json:"apiVersion"`
	// The Kind of the source; e.g., `GitRepository`
	Kind string `json:"kind"`
	// The name of the source.
	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StackList contains a list of Stack
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stack{}, &StackList{})
}
