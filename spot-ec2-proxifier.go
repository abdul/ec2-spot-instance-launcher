package main;

import(
	"fmt"
	"errors"
	"os/exec"
	"github.com/msbranco/goconfig"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
);

func main() {
	config, configError := goconfig.ReadConfigFile("spot-ec2-proxifier.ini");

	if configError != nil {
		fmt.Println(configError);
		return;
	}

	// Setup our AWS credentials
	ak, _ := config.GetString("IAM", "access")
	sk, _ := config.GetString("IAM", "secret")
	rn, _ := config.GetString("EC2", "region")

	e, err := setupEC2(ak, sk, rn);
	if err != nil {
		fmt.Println("Something wrong in EC2 config:", err);
		return;
	}

	// Get all current running instances 
	rep, err := e.Instances(nil, nil)
	if err != nil {
		fmt.Println("Couldn't list instances", err);
		return;
	}

	// Find whether any are ours
	tag, _ := config.GetString("EC2", "tag");
	inst, found := findTaggedInstance(rep.Reservations, tag);

	if found {
		fmt.Println("Found tagged instance:", inst.DNSName);

		// Launch plink on that instance
		bindPort, _ := config.GetInt64("Proxy", "bind_port");
		keyPath, _ := config.GetString("Proxy", "key_file");
		username, _ := config.GetString("EC2", "username");

		if username == "" {
			username = "ec2-user";
		}

		connectPlink(username, inst.DNSName, keyPath, int(bindPort));

	} else {
		fmt.Println("No instance up, creating a new one");
	}
}

func setupEC2(accessKey string, secretKey string, regionName string) (ec2.EC2, error)  {
	e := ec2.EC2{};

	auth := aws.Auth { accessKey, secretKey };

	region := aws.Regions[regionName]
	if region.Name != regionName {
		return e, errors.New("Region is invalid");
	}

	e.Auth = auth;
	e.Region = region;

	return e, nil;
}

func findTaggedInstance(reservations []ec2.Reservation, search_tag string) (ec2.Instance, bool) {
	for _, res := range reservations {
		for _, instance := range res.Instances {
			for _, tag := range instance.Tags {
				if tag.Value == search_tag {
					return instance, true;
				}
			}
		}
	}
	return ec2.Instance{}, false;
}

func connectPlink(username string, host string, keyFile string, bindPort int) {
	cmd := exec.Command("plink", "-N", "-D", fmt.Sprintf("%d", bindPort), "-i", keyFile, fmt.Sprintf("%s@%s", username, host));

	err := cmd.Start();
	if err != nil {
		fmt.Println("Error:", err);
		return;
	}	

	fmt.Println("Use ^c to kill plink (and this application).");
	cmd.Wait();
}