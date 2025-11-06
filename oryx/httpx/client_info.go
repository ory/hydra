// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

type GeoLocation struct {
	City      string
	Region    string
	Country   string
	Latitude  *float64
	Longitude *float64
}

func GetClientIPAddressesWithoutInternalIPs(ipAddresses []string) (string, error) {
	var res string

	for i := len(ipAddresses) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(ipAddresses[i])

		if !net.ParseIP(ip).IsPrivate() {
			res = ip
			break
		}
	}

	return res, nil
}

func ClientIP(r *http.Request) string {
	if trueClientIP := r.Header.Get("True-Client-IP"); trueClientIP != "" {
		return trueClientIP
	} else if cfConnectingIP := r.Header.Get("Cf-Connecting-IP"); cfConnectingIP != "" {
		return cfConnectingIP
	} else if realClientIP := r.Header.Get("X-Real-IP"); realClientIP != "" {
		return realClientIP
	} else if forwardedIP := r.Header.Get("X-Forwarded-For"); forwardedIP != "" {
		ip, _ := GetClientIPAddressesWithoutInternalIPs(strings.Split(forwardedIP, ","))
		return ip
	} else {
		return r.RemoteAddr
	}
}

func parseFloatHeaderValue(headerValue string) *float64 {
	if headerValue == "" {
		return nil
	}

	val, err := strconv.ParseFloat(headerValue, 64)
	if err != nil {
		return nil
	}

	return &val
}

func ClientGeoLocation(r *http.Request) *GeoLocation {
	return &GeoLocation{
		City:      r.Header.Get("Cf-Ipcity"),
		Region:    r.Header.Get("Cf-Region-Code"),
		Country:   r.Header.Get("Cf-Ipcountry"),
		Longitude: parseFloatHeaderValue(r.Header.Get("Cf-Iplongitude")),
		Latitude:  parseFloatHeaderValue(r.Header.Get("Cf-Iplatitude")),
	}
}
