package serve

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"tailproxy/src/config"
	"tailproxy/src/ts"
	"time"
)

func makeProxy() http.Handler {
	var start time.Time
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			start = time.Now()
			fmt.Printf("tailproxy: %v %v %v\n", r.In.RemoteAddr, r.In.Method, r.In.URL)
			r.SetXForwarded()
			r.SetURL(config.Target)
			r.Out.Host = r.In.Host

			who, err := ts.WhoIs(r.In)
			for k := range r.Out.Header {
				if strings.HasPrefix(k, "X-Tailscale-") {
					r.Out.Header.Del(k)
				}
			}
			if err != nil {
				fmt.Printf("error getting whois: %v\n", err)
				r.Out.Header.Set("X-Tailscale-WhoIs", "error")
			} else {
				r.Out.Header.Set("X-Tailscale-WhoIs", "ok")
				r.Out.Header.Set("X-Tailscale-User", who.UserProfile.ID.String())
				r.Out.Header.Set("X-Tailscale-User-LoginName", who.UserProfile.LoginName)
				r.Out.Header.Set("X-Tailscale-User-DisplayName", who.UserProfile.DisplayName)
				if who.UserProfile.ProfilePicURL != "" {
					r.Out.Header.Set("X-Tailscale-User-ProfilePicURL", who.UserProfile.ProfilePicURL)
				}

				if len(who.Caps) > 0 {
					r.Out.Header.Set("X-Tailscale-Caps", strings.Join(who.Caps, ", "))
				}

				r.Out.Header.Set("X-Tailscale-Node", who.Node.ID.String())
				r.Out.Header.Set("X-Tailscale-Node-Name", who.Node.ComputedName)
				if len(who.Node.Capabilities) > 0 {
					r.Out.Header.Set("X-Tailscale-Node-Caps", strings.Join(who.Node.Capabilities, ", "))
				}
				if len(who.Node.Tags) > 0 {
					r.Out.Header.Set("X-Tailscale-Node-Tags", strings.Join(who.Node.Tags, ", "))
				}
				data, err := who.Node.Hostinfo.MarshalJSON()
				if err == nil {
					r.Out.Header.Set("X-Tailscale-Hostinfo", string(data))
				}
			}
		},
		ModifyResponse: func(r *http.Response) error {
			fmt.Printf("tailproxy: %v %v %v %v %v\n", r.Request.RemoteAddr, r.Request.Method, r.Request.URL, r.StatusCode, time.Since(start))
			return nil
		},
	}
}
