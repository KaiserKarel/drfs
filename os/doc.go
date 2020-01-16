// Package os implements file related functions similar to the standard library os package. It
// attempts to be a drop in replacement using drfs to reduce boilerplate code. Before usage a
// drfs client is initiated by looking for a secret either in the directory specified by
// DRFS_APPLICATION_CREDENTIALS or a specific directory using GOOGLE_APPLICATION_CREDENTIALS.
package os
