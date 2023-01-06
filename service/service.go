package main

import (
	"context"
	"fmt"
	"github.com/MarcGrol/patientregistration/lib/api/datastorer"
	"github.com/MarcGrol/patientregistration/lib/api/emailsender"
	"github.com/MarcGrol/patientregistration/lib/api/pingenerator"
	"github.com/MarcGrol/patientregistration/regprotobuf"
	"github.com/google/uuid"
	"strconv"
)

type RegistrationService struct {
	patientStore datastorer.PatientStorer
	emailSender  emailsender.EmailSender
	pinGenerator pingenerator.PinGenerator
	regprotobuf.UnimplementedRegistrationServiceServer
}

func NewRegistrationService(
	patientStore datastorer.PatientStorer,
	emailSender emailsender.EmailSender,
	pinGenerator pingenerator.PinGenerator,
) *RegistrationService {
	return &RegistrationService{
		patientStore: patientStore,
		emailSender:  emailSender,
		pinGenerator: pinGenerator,
	}
}

func (rs *RegistrationService) RegisterPatient(ctx context.Context, req *regprotobuf.RegisterPatientRequest) (*regprotobuf.RegisterPatientResponse, error) {
	email := req.GetPatient().GetContact().GetEmailAddress()

	// Validate email is already registered
	patient, _, err := rs.patientStore.GetPatientOnEmail(email)
	if err == nil && patient.UID != "" {
		fmt.Println("Email is already registered: ", email)
		return nil, err
	}

	// generate pin
	pin, _ := rs.pinGenerator.GeneratePin()

	//send pin to email
	err = rs.emailSender.SendEmail(email, "validate", "This is your pin "+strconv.Itoa(pin))
	if err != nil {
		fmt.Println("Cannot send PIN to: ", email)
		return nil, err
	}

	// store user on database as pending
	patient = datastorer.Patient{
		UID:      uuid.New().String(),
		BSN:      req.Patient.BSN,
		FullName: req.Patient.FullName,
		Address: datastorer.StreetAddress{
			PostalCode:  req.Patient.Address.PostalCode,
			HouseNumber: int(req.Patient.Address.HouseNumber),
		},
		Contact: datastorer.Contact{
			EmailAddress: req.Patient.Contact.EmailAddress,
		},
		RegistrationPin:    pin,
		FailedPinCount:     0,
		RegistrationStatus: datastorer.Pending,
	}
	err = rs.patientStore.PutPatientOnUid(patient)
	if err != nil {
		fmt.Println("Cannot register user: ", email)
		return nil, err
	}

	return &regprotobuf.RegisterPatientResponse{
		PatientUid: patient.UID,
	}, nil
}

// TODO add CompletePatientRegistration
