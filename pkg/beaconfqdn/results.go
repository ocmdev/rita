package beaconfqdn

import (
	"github.com/activecm/rita/resources"
	"github.com/globalsign/mgo/bson"
)

//Results finds beacons FQDN in the database greater than a given cutoffScore
func Results(res *resources.Resources, cutoffScore float64) ([]Result, error) {
	ssn := res.DB.Session.Copy()
	defer ssn.Close()

	var beaconsFQDN []Result

	beaconFQDNQuery := bson.M{"score": bson.M{"$gt": cutoffScore}}

	err := ssn.DB(res.DB.GetSelectedDB()).C(res.Config.T.BeaconFQDN.BeaconFQDNTable).Find(beaconFQDNQuery).Sort("-score").All(&beaconsFQDN)

	return beaconsFQDN, err
}
