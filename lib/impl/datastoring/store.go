package datastoring

import (
	"github.com/MarcGrol/patientregistration/lib/api/datastorer"
	"sync"
)

type inMemoryPatientStore struct {
	sync.Mutex
	patients map[string]datastorer.Patient
}

func New() datastorer.PatientStorer {
	return &inMemoryPatientStore{
		patients: map[string]datastorer.Patient{},
	}
}
func (ps *inMemoryPatientStore) PutPatientOnUid(patient datastorer.Patient) error {
	ps.Lock()
	defer ps.Unlock()

	ps.patients[patient.UID] = patient

	return nil
}

func (ps *inMemoryPatientStore) GetPatientOnUid(patientUID string) (datastorer.Patient, bool, error) {
	ps.Lock()
	defer ps.Unlock()

	patient, found := ps.patients[patientUID]

	return patient, found, nil
}

func (ps *inMemoryPatientStore) GetPatientOnEmail(patientEmail string) (datastorer.Patient, bool, error) {
	ps.Lock()
	defer ps.Unlock()

	for patientUID := range ps.patients {
		patient := ps.patients[patientUID]
		if patientEmail == ps.patients[patientUID].Contact.EmailAddress {
			return patient, true, nil
		}
	}

	return datastorer.Patient{}, false, nil
}
