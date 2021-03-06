%define  debug_package %{nil}

Name:    cwlogs-tee
Version: 0.1.4
Release: 1%{?dist}
Summary: cwlogs-tee is a tee command for CloudWatch Logs.

Group:   Development/Tools
License: MIT License
URL:     https://github.com/winebarrel/cwlogs-tee
Source0: %{name}.tar.gz

%description
cwlogs-tee is a tee command for CloudWatch Logs.

%prep
%setup -q -n src

%build
make
make test

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
install -m 755 cwlogs-tee %{buildroot}/usr/bin/

%files
%defattr(755,root,root,-)
/usr/bin/cwlogs-tee
