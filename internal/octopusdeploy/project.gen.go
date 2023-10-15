// This code is generated, DO NOT EDIT
package octopusdeploy

type Project struct {
	AutoCreateRelease               bool          `json:"AutoCreateRelease"`
	AutoDeployReleaseOverrides      []interface{} `json:"AutoDeployReleaseOverrides"`
	ClonedFromProjectID             interface{}   `json:"ClonedFromProjectId"`
	DefaultGuidedFailureMode        string        `json:"DefaultGuidedFailureMode"`
	DefaultToSkipIfAlreadyInstalled bool          `json:"DefaultToSkipIfAlreadyInstalled"`
	DeploymentChangesTemplate       interface{}   `json:"DeploymentChangesTemplate"`
	DeploymentProcessID             string        `json:"DeploymentProcessId"`
	Description                     string        `json:"Description"`
	DiscreteChannelRelease          bool          `json:"DiscreteChannelRelease"`
	ExtensionSettings               []interface{} `json:"ExtensionSettings"`
	ForcePackageDownload            bool          `json:"ForcePackageDownload"`
	Icon                            interface{}   `json:"Icon"`
	ID                              string        `json:"Id"`
	IncludedLibraryVariableSetIds   []interface{} `json:"IncludedLibraryVariableSetIds"`
	IsDisabled                      bool          `json:"IsDisabled"`
	IsVersionControlled             bool          `json:"IsVersionControlled"`
	LifecycleID                     string        `json:"LifecycleId"`
	Links                           struct {
		Channels                             string `json:"Channels"`
		ConvertToGit                         string `json:"ConvertToGit"`
		ConvertToVcs                         string `json:"ConvertToVcs"`
		DeploymentProcess                    string `json:"DeploymentProcess"`
		DeploymentSettings                   string `json:"DeploymentSettings"`
		GitCompatibilityReport               string `json:"GitCompatibilityReport"`
		GitConnectionTest                    string `json:"GitConnectionTest"`
		InsightsMetrics                      string `json:"InsightsMetrics"`
		Logo                                 string `json:"Logo"`
		Metadata                             string `json:"Metadata"`
		OrderChannels                        string `json:"OrderChannels"`
		Progression                          string `json:"Progression"`
		Releases                             string `json:"Releases"`
		RunbookSnapshots                     string `json:"RunbookSnapshots"`
		RunbookTaskRunDashboardItemsTemplate string `json:"RunbookTaskRunDashboardItemsTemplate"`
		Runbooks                             string `json:"Runbooks"`
		ScheduledTriggers                    string `json:"ScheduledTriggers"`
		Self                                 string `json:"Self"`
		Summary                              string `json:"Summary"`
		Triggers                             string `json:"Triggers"`
		Variables                            string `json:"Variables"`
		Web                                  string `json:"Web"`
	} `json:"Links"`
	Name                string `json:"Name"`
	PersistenceSettings struct {
		Type string `json:"Type"`
	} `json:"PersistenceSettings"`
	ProjectConnectivityPolicy struct {
		AllowDeploymentsToNoTargets bool          `json:"AllowDeploymentsToNoTargets"`
		ExcludeUnhealthyTargets     bool          `json:"ExcludeUnhealthyTargets"`
		SkipMachineBehavior         string        `json:"SkipMachineBehavior"`
		TargetRoles                 []interface{} `json:"TargetRoles"`
	} `json:"ProjectConnectivityPolicy"`
	ProjectGroupID          string `json:"ProjectGroupId"`
	ReleaseCreationStrategy struct {
		ChannelID                    interface{} `json:"ChannelId"`
		ReleaseCreationPackage       interface{} `json:"ReleaseCreationPackage"`
		ReleaseCreationPackageStepID interface{} `json:"ReleaseCreationPackageStepId"`
	} `json:"ReleaseCreationStrategy"`
	ReleaseNotesTemplate   interface{}   `json:"ReleaseNotesTemplate"`
	Slug                   string        `json:"Slug"`
	SpaceID                string        `json:"SpaceId"`
	Templates              []interface{} `json:"Templates"`
	TenantedDeploymentMode string        `json:"TenantedDeploymentMode"`
	VariableSetID          string        `json:"VariableSetId"`
	VersioningStrategy     struct {
		DonorPackage interface{} `json:"DonorPackage"`
		Template     string      `json:"Template"`
	} `json:"VersioningStrategy"`
}
