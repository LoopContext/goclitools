package goclitools

// DependencyHomebrew dependency homebrew
func DependencyHomebrew() Dependency {
	return Dependency{
		Name:     "Homebrew",
		CheckCmd: "which brew",
		InstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""},
		},
	}
}

// DependencyDocker docker dependency
func DependencyDocker() Dependency {
	return Dependency{
		Name:         "Docker",
		CheckCmd:     "which docker",
		Dependencies: []Dependency{DependencyHomebrew()},
		InstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "brew cask install docker"},
			DependencyScriptString{Fn: "open /Applications/Docker.app"},
		},
		UninstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "brew cask uninstall docker"},
		},
	}
}

// DependencyGit git dependency
func DependencyGit() Dependency {
	return Dependency{
		Name:               "git",
		CheckCmd:           "git --version",
		CheckCmdValidation: "(?m)git version (\\d+\\.)?(\\d+\\.)?(\\*|\\d+).*$",
	}
}

// DependencyXcodebuild Xcode build
func DependencyXcodebuild() Dependency {
	return Dependency{
		Name:               "xcodebuild",
		CheckCmd:           "xcodebuild -version",
		CheckCmdValidation: "(?m)Xcode (\\d+\\.)?(\\d+\\.)?(\\*|\\d+)\\s+Build version .+$",
	}
}

// DependencyFastlane fastlane dependency
func DependencyFastlane() Dependency {
	return Dependency{
		Name:               "fastlane",
		CheckCmd:           "fastlane -version",
		CheckCmdValidation: "(?m)fastlane (\\d+\\.)?(\\d+\\.)?(\\*|\\d+)$",
		Dependencies:       []Dependency{DependencyHomebrew()},
		InstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "brew cask install fastlane"},
			DependencyScriptString{Fn: "export PATH=\"$HOME/.fastlane/bin:$PATH\""},
			DependencyScriptString{Fn: "echo \"export PATH=\\\"$HOME/.fastlane/bin:$PATH\\\"\" >> ~/.bash_profile"},
		},
		UninstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "brew cask uninstall fastlane"},
		},
	}
}

// DependencyFastlaneMatch fastlane dependency
func DependencyFastlaneMatch() Dependency {
	return Dependency{
		Name:               "fastlane match",
		CheckCmd:           "fastlane match -version",
		CheckCmdValidation: "(?m)match (\\d+\\.)?(\\d+\\.)?(\\*|\\d+)$",
		Dependencies:       []Dependency{DependencyFastlane()},
		InstallScripts:     []DependencyScript{},
		UninstallScripts: []DependencyScript{
			DependencyScriptString{Fn: "echo \"you can only uninstall whole fastlane\""},
		},
	}
}
