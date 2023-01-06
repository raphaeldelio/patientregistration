package main

import (
	"context"
	"github.com/MarcGrol/patientregistration/lib/api/datastorer"
	"github.com/MarcGrol/patientregistration/lib/api/emailsender"
	"github.com/MarcGrol/patientregistration/lib/api/pingenerator"
	"github.com/MarcGrol/patientregistration/regprotobuf"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"testing"
)

func newRegistrationService(t *testing.T) *RegistrationService {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	return NewRegistrationService(
		datastorer.NewMockPatientStorer(ctrl),
		emailsender.NewMockEmailSender(ctrl),
		pingenerator.NewMockPinGenerator(ctrl),
	)
}

func TestRegistrationWithEmailSuccess(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataStorer := datastorer.NewMockPatientStorer(ctrl)
	mockEmailSender := emailsender.NewMockEmailSender(ctrl)
	mockPinGenerator := pingenerator.NewMockPinGenerator(ctrl)

	service := NewRegistrationService(
		mockDataStorer,
		mockEmailSender,
		mockPinGenerator,
	)

	ctx := context.Background()

	// MOCKS
	mockDataStorer.EXPECT().GetPatientOnEmail("john@xebia.com").Return(datastorer.Patient{}, false, nil)
	mockPinGenerator.EXPECT().GeneratePin().Return(1111, nil)
	mockEmailSender.EXPECT().SendEmail("john@xebia.com", gomock.Any(), gomock.Any()).Return(nil)
	mockDataStorer.EXPECT().PutPatientOnUid(gomock.Any()).Return(nil)

	// WHEN
	request := regprotobuf.RegisterPatientRequest{
		Patient: &regprotobuf.Patient{
			BSN:      "",
			FullName: "",
			Address: &regprotobuf.Address{
				PostalCode:  "",
				HouseNumber: 0,
			},
			Contact: &regprotobuf.Contact{
				EmailAddress: "john@xebia.com",
			},
			Status: 0,
		},
	}

	response, err := service.RegisterPatient(ctx, &request)

	// THEN
	assert.NoError(t, err)
	assert.NotEmpty(t, response.PatientUid)
}

func TestRegistrationInvalidInput(t *testing.T) {
	// TODO
}

func TestRegistrationSendEmailError(t *testing.T) {
	// TODO
}
