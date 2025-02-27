/*
 * Siren.
 *
 * Documentation of our Siren API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

type Annotations struct {
	MetricName string `json:"metric_name,omitempty"`
	MetricValue string `json:"metric_value,omitempty"`
	Resource string `json:"resource,omitempty"`
	Template string `json:"template,omitempty"`
}
