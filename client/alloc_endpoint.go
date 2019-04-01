package client

import (
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/nomad/acl"
	cstructs "github.com/hashicorp/nomad/client/structs"
	nstructs "github.com/hashicorp/nomad/nomad/structs"
)

// Allocations endpoint is used for interacting with client allocations
type Allocations struct {
	c *Client
}

// GarbageCollectAll is used to garbage collect all allocations on a client.
func (a *Allocations) GarbageCollectAll(args *nstructs.NodeSpecificRequest, reply *nstructs.GenericResponse) error {
	defer metrics.MeasureSince([]string{"client", "allocations", "garbage_collect_all"}, time.Now())

	// Check node write permissions
	if aclObj, err := a.c.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowNodeWrite() {
		return nstructs.ErrPermissionDenied
	}

	a.c.CollectAllAllocs()
	return nil
}

// GarbageCollect is used to garbage collect an allocation on a client.
func (a *Allocations) GarbageCollect(args *nstructs.AllocSpecificRequest, reply *nstructs.GenericResponse) error {
	defer metrics.MeasureSince([]string{"client", "allocations", "garbage_collect"}, time.Now())

	// Check submit job permissions
	if aclObj, err := a.c.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowNsOp(args.Namespace, acl.NamespaceCapabilitySubmitJob) {
		return nstructs.ErrPermissionDenied
	}

	if !a.c.CollectAllocation(args.AllocID) {
		// Could not find alloc
		return nstructs.NewErrUnknownAllocation(args.AllocID)
	}

	return nil
}

func (a *Allocations) Signal(args *nstructs.AllocSignalRequest, reply *nstructs.GenericResponse) error {
	defer metrics.MeasureSince([]string{"client", "allocations", "signal"}, time.Now())

	// Check submit job permissions
	if aclObj, err := a.c.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowNsOp(args.Namespace, acl.NamespaceCapabilityAllocLifecycle) {
		return nstructs.ErrPermissionDenied
	}

	return a.c.SignalAllocation(args.AllocID, args.Task, args.Signal)
}

// Restart is used to trigger a restart of an allocation or a subtask on a client.
func (a *Allocations) Restart(args *nstructs.AllocRestartRequest, reply *nstructs.GenericResponse) error {
	defer metrics.MeasureSince([]string{"client", "allocations", "restart"}, time.Now())

	if aclObj, err := a.c.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowNsOp(args.Namespace, acl.NamespaceCapabilityAllocLifecycle) {
		return nstructs.ErrPermissionDenied
	}

	return a.c.RestartAllocation(args.AllocID, args.TaskName)
}

// Stats is used to collect allocation statistics
func (a *Allocations) Stats(args *cstructs.AllocStatsRequest, reply *cstructs.AllocStatsResponse) error {
	defer metrics.MeasureSince([]string{"client", "allocations", "stats"}, time.Now())

	// Check read job permissions
	if aclObj, err := a.c.ResolveToken(args.AuthToken); err != nil {
		return err
	} else if aclObj != nil && !aclObj.AllowNsOp(args.Namespace, acl.NamespaceCapabilityReadJob) {
		return nstructs.ErrPermissionDenied
	}

	clientStats := a.c.StatsReporter()
	aStats, err := clientStats.GetAllocStats(args.AllocID)
	if err != nil {
		return err
	}

	stats, err := aStats.LatestAllocStats(args.Task)
	if err != nil {
		return err
	}

	reply.Stats = stats
	return nil
}
