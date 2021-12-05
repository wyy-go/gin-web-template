package env

type Deploy int8

const (
	DeployUnknown Deploy = iota
	DeployLocal
	DeployDev
	DeployTest
	DeployUat
	DeployProd
)

func (d Deploy) String() string {
	switch d {
	case DeployLocal:
		return "local"
	case DeployDev:
		return "dev"
	case DeployTest:
		return "test"
	case DeployUat:
		return "uat"
	case DeployProd:
		return "prod"
	}
	return ""
}

func ToDeploy(s string) Deploy {
	switch s {
	case DeployLocal.String():
		return DeployLocal
	case DeployDev.String():
		return DeployDev
	case DeployTest.String():
		return DeployTest
	case DeployUat.String():
		return DeployUat
	case DeployProd.String():
		return DeployProd
	}
	return DeployUnknown
}

var deploy = DeployDev

func GetDeploy() Deploy {
	return deploy
}

func SetDeploy(d Deploy) {
	deploy = d
}

func IsDeployLocal() bool {
	return deploy == DeployLocal
}

func IsDeployDev() bool {
	return deploy == DeployDev
}

func IsDeployTest() bool {
	return deploy == DeployTest
}

func IsDeployUat() bool {
	return deploy == DeployUat
}

func IsDeployProd() bool {
	return deploy == DeployProd
}

func IsDeployDebug() bool {
	return IsDeployLocal() || IsDeployDev() || IsDeployTest()
}

func IsDeployRelease() bool {
	return IsDeployUat() || IsDeployProd()
}
