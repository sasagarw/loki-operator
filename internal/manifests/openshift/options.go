package openshift

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
)

// Options is the set of internal template options for rendering
// the lokistack-gateway tenants configuration file when mode openshift-logging.
type Options struct {
	BuildOpts      BuildOptions
	Authentication []AuthenticationSpec
	Authorization  AuthorizationSpec
}

// AuthenticationSpec describes the authentication specification
// for a single tenant to authenticate it's subjects through OpenShift Auth.
type AuthenticationSpec struct {
	TenantName     string
	TenantID       string
	ServiceAccount string
	RedirectURL    string
	CookieSecret   string
}

// AuthorizationSpec describes the authorization specification
// for all tenants to authorize access for it's subjects through the
// opa-openshift sidecar.
type AuthorizationSpec struct {
	OPAUrl string
}

// BuildOptions represents the set of options required to build
// extra lokistack gateway k8s objects (e.g. ServiceAccount, Route, RBAC)
// on openshift.
type BuildOptions struct {
	LokiStackName                   string
	GatewayName                     string
	GatewayNamespace                string
	GatewaySvcName                  string
	GatewaySvcTargetPort            string
	Labels                          map[string]string
	EnableCertificateSigningService bool
}

// NewOptions returns an openshift options struct.
func NewOptions(
	stackName string,
	gwName, gwNamespace, gwBaseDomain, gwSvcName, gwPortName string,
	gwLabels map[string]string,
	enableCertSigningService bool,
) Options {
	host := ingressHost(stackName, gwNamespace, gwBaseDomain)

	var authn []AuthenticationSpec
	for _, name := range defaultTenants {
		s := sha1.New()
		s.Write([]byte(name))
		tenantID := s.Sum(nil)
		authn = append(authn, AuthenticationSpec{
			TenantName:     name,
			TenantID:       fmt.Sprintf("%x", tenantID),
			ServiceAccount: gwName,
			RedirectURL:    fmt.Sprintf("http://%s/openshift/%s/callback", host, name),
			CookieSecret:   newCookieSecret(),
		})
	}

	return Options{
		BuildOpts: BuildOptions{
			LokiStackName:                   stackName,
			GatewayName:                     gwName,
			GatewayNamespace:                gwNamespace,
			GatewaySvcName:                  gwSvcName,
			GatewaySvcTargetPort:            gwPortName,
			Labels:                          gwLabels,
			EnableCertificateSigningService: enableCertSigningService,
		},
		Authentication: authn,
		Authorization: AuthorizationSpec{
			OPAUrl: fmt.Sprintf("http://localhost:%d/v1/data/%s/allow", GatewayOPAHTTPPort, opaDefaultPackage),
		},
	}
}

func newCookieSecret() string {
	b := make([]rune, cookieSecretLength)
	for i := range b {
		b[i] = allowedRunes[rand.Intn(len(allowedRunes))]
	}

	return string(b)
}
