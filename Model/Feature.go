package model

type Feature struct {
	FeatureID          int    `json:"featureID"`
	FeatureName        string `json:"featureName"`
	FeatureType        string `json:"featureType"`
	FeatureDescription string `json:"featureDescription"`
}
