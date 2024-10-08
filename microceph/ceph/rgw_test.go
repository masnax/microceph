package ceph

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/microceph/microceph/tests"

	"github.com/canonical/microceph/microceph/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type rgwSuite struct {
	tests.BaseSuite
	TestStateInterface *mocks.StateInterface
}

func TestRGW(t *testing.T) {
	suite.Run(t, new(rgwSuite))
}

// Expect: run ceph auth
func addRGWEnableExpectations(r *mocks.Runner) {
	// add keyring expectation
	r.On("RunCommand", tests.CmdAny("ceph", 9)...).Return("ok", nil).Once()
	// start service expectation
	r.On("RunCommand", []interface{}{
		"snapctl", "start", "microceph.rgw", "--enable",
	}...).Return("ok", nil).Once()
}

// Expect: run snapctl service stop
func addStopRGWExpectations(s *rgwSuite, r *mocks.Runner) {
	u := api.NewURL()

	state := &mocks.MockState{
		URL:         u,
		ClusterName: "foohost",
	}

	s.TestStateInterface.On("ClusterState").Return(state)
	r.On("RunCommand", tests.CmdAny("snapctl", 3)...).Return("ok", nil).Once()
}

// Set up test suite
func (s *rgwSuite) SetupTest() {
	s.BaseSuite.SetupTest()
	s.CopyCephConfigs()

	s.TestStateInterface = mocks.NewStateInterface(s.T())
}

// Test enabling RGW
func (s *rgwSuite) TestEnableRGW() {
	r := mocks.NewRunner(s.T())

	addRGWEnableExpectations(r)

	processExec = r

	err := EnableRGW(s.TestStateInterface, 80, []string{"10.1.1.1", "10.2.2.2"})

	assert.NoError(s.T(), err)

	// check that the radosgw.conf file contains expected values
	conf := s.ReadCephConfig("radosgw.conf")
	assert.Contains(s.T(), conf, "rgw frontends = beast port=80")
	assert.Contains(s.T(), conf, "mon host = 10.1.1.1,10.2.2.2")
}

func (s *rgwSuite) TestDisableRGW() {
	r := mocks.NewRunner(s.T())

	addStopRGWExpectations(s, r)

	processExec = r

	err := DisableRGW(context.Background(), s.TestStateInterface)

	// we expect a missing database error
	assert.EqualError(s.T(), err, "no server certificate")

	// check that the radosgw.conf file is absent
	_, err = os.Stat(filepath.Join(s.Tmp, "SNAP_DATA", "conf", "radosgw.conf"))
	assert.True(s.T(), os.IsNotExist(err))

	// check that the keyring file is absent
	_, err = os.Stat(filepath.Join(s.Tmp, "SNAP_COMMON", "data", "radosgw", "ceph-radosgw.gateway", "keyring"))
	assert.True(s.T(), os.IsNotExist(err))
}
