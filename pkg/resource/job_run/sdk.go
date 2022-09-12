// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package job_run

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/emrcontainers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/emrcontainers-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.EMRContainers{}
	_ = &svcapitypes.JobRun{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeJobRunOutput
	resp, err = rm.sdkapi.DescribeJobRunWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeJobRun", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "UNKNOWN" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.JobRun.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.JobRun.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.JobRun.ExecutionRoleArn != nil {
		ko.Spec.ExecutionRoleARN = resp.JobRun.ExecutionRoleArn
	} else {
		ko.Spec.ExecutionRoleARN = nil
	}
	if resp.JobRun.Id != nil {
		ko.Status.ID = resp.JobRun.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.JobRun.JobDriver != nil {
		f8 := &svcapitypes.JobDriver{}
		if resp.JobRun.JobDriver.SparkSubmitJobDriver != nil {
			f8f0 := &svcapitypes.SparkSubmitJobDriver{}
			if resp.JobRun.JobDriver.SparkSubmitJobDriver.EntryPoint != nil {
				f8f0.EntryPoint = resp.JobRun.JobDriver.SparkSubmitJobDriver.EntryPoint
			}
			if resp.JobRun.JobDriver.SparkSubmitJobDriver.EntryPointArguments != nil {
				f8f0f1 := []*string{}
				for _, f8f0f1iter := range resp.JobRun.JobDriver.SparkSubmitJobDriver.EntryPointArguments {
					var f8f0f1elem string
					f8f0f1elem = *f8f0f1iter
					f8f0f1 = append(f8f0f1, &f8f0f1elem)
				}
				f8f0.EntryPointArguments = f8f0f1
			}
			if resp.JobRun.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters != nil {
				f8f0.SparkSubmitParameters = resp.JobRun.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters
			}
			f8.SparkSubmitJobDriver = f8f0
		}
		ko.Spec.JobDriver = f8
	} else {
		ko.Spec.JobDriver = nil
	}
	if resp.JobRun.Name != nil {
		ko.Spec.Name = resp.JobRun.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.JobRun.ReleaseLabel != nil {
		ko.Spec.ReleaseLabel = resp.JobRun.ReleaseLabel
	} else {
		ko.Spec.ReleaseLabel = nil
	}
	if resp.JobRun.Tags != nil {
		f13 := map[string]*string{}
		for f13key, f13valiter := range resp.JobRun.Tags {
			var f13val string
			f13val = *f13valiter
			f13[f13key] = &f13val
		}
		ko.Spec.Tags = f13
	} else {
		ko.Spec.Tags = nil
	}
	if resp.JobRun.VirtualClusterId != nil {
		ko.Spec.VirtualClusterID = resp.JobRun.VirtualClusterId
	} else {
		ko.Spec.VirtualClusterID = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Status.ID == nil || r.ko.Spec.VirtualClusterID == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeJobRunInput, error) {
	res := &svcsdk.DescribeJobRunInput{}

	if r.ko.Status.ID != nil {
		res.SetId(*r.ko.Status.ID)
	}
	if r.ko.Spec.VirtualClusterID != nil {
		res.SetVirtualClusterId(*r.ko.Spec.VirtualClusterID)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.StartJobRunOutput
	_ = resp
	resp, err = rm.sdkapi.StartJobRunWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "StartJobRun", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Id != nil {
		ko.Status.ID = resp.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.Name != nil {
		ko.Spec.Name = resp.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.VirtualClusterId != nil {
		ko.Spec.VirtualClusterID = resp.VirtualClusterId
	} else {
		ko.Spec.VirtualClusterID = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.StartJobRunInput, error) {
	res := &svcsdk.StartJobRunInput{}

	if r.ko.Spec.ExecutionRoleARN != nil {
		res.SetExecutionRoleArn(*r.ko.Spec.ExecutionRoleARN)
	}
	if r.ko.Spec.JobDriver != nil {
		f1 := &svcsdk.JobDriver{}
		if r.ko.Spec.JobDriver.SparkSubmitJobDriver != nil {
			f1f0 := &svcsdk.SparkSubmitJobDriver{}
			if r.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint != nil {
				f1f0.SetEntryPoint(*r.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint)
			}
			if r.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments != nil {
				f1f0f1 := []*string{}
				for _, f1f0f1iter := range r.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments {
					var f1f0f1elem string
					f1f0f1elem = *f1f0f1iter
					f1f0f1 = append(f1f0f1, &f1f0f1elem)
				}
				f1f0.SetEntryPointArguments(f1f0f1)
			}
			if r.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters != nil {
				f1f0.SetSparkSubmitParameters(*r.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters)
			}
			f1.SetSparkSubmitJobDriver(f1f0)
		}
		res.SetJobDriver(f1)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.ReleaseLabel != nil {
		res.SetReleaseLabel(*r.ko.Spec.ReleaseLabel)
	}
	if r.ko.Spec.Tags != nil {
		f4 := map[string]*string{}
		for f4key, f4valiter := range r.ko.Spec.Tags {
			var f4val string
			f4val = *f4valiter
			f4[f4key] = &f4val
		}
		res.SetTags(f4)
	}
	if r.ko.Spec.VirtualClusterID != nil {
		res.SetVirtualClusterId(*r.ko.Spec.VirtualClusterID)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.CancelJobRunOutput
	_ = resp
	resp, err = rm.sdkapi.CancelJobRunWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "CancelJobRun", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.CancelJobRunInput, error) {
	res := &svcsdk.CancelJobRunInput{}

	if r.ko.Status.ID != nil {
		res.SetId(*r.ko.Status.ID)
	}
	if r.ko.Spec.VirtualClusterID != nil {
		res.SetVirtualClusterId(*r.ko.Spec.VirtualClusterID)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.JobRun,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "ValidationException":
		return true
	default:
		return false
	}
}

// getImmutableFieldChanges returns list of immutable fields from the
func (rm *resourceManager) getImmutableFieldChanges(
	delta *ackcompare.Delta,
) []string {
	var fields []string
	if delta.DifferentAt("Spec.ExecutionRoleARN") {
		fields = append(fields, "ExecutionRoleARN")
	}
	if delta.DifferentAt("Spec.JobDriver") {
		fields = append(fields, "JobDriver")
	}
	if delta.DifferentAt("Spec.Name") {
		fields = append(fields, "Name")
	}
	if delta.DifferentAt("Spec.ReleaseLabel") {
		fields = append(fields, "ReleaseLabel")
	}
	if delta.DifferentAt("Spec.VirtualClusterId") {
		fields = append(fields, "VirtualClusterId")
	}

	return fields
}
