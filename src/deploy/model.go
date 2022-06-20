package deploy

type dep interface {
	Source()
	Curl()
	Wget()
	Supervisor()
	Jq()
	Wireguard()
}

type ItsOUin struct {
	OS   string
	Type string
	name []string
}

const (
	aptStr = `
deb $apt_url $i.name main restricted universe multiverse
deb-src $apt_url $i.name main restricted universe multiverse
deb $apt_url $i.name-security main restricted universe multiverse
deb-src $apt_url $i.name-security main restricted universe multiverse
deb $apt_url $i.name-updates main restricted universe multiverse
deb-src $apt_url $i.name-updates main restricted universe multiverse
deb $apt_url $i.name-proposed main restricted universe multiverse
deb-src $apt_url $i.name-proposed main restricted universe multiverse
deb $apt_url $i.name-backports main restricted universe multiverse
deb-src $apt_url $i.name-backports main restricted universe multiverse`
	aptUrl = `http://archive.ubuntu.com/ubuntu/`
)
