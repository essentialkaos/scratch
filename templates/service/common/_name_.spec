################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define debug_package  %{nil}

################################################################################

%define _logdir  %{_localstatedir}/log

################################################################################

Summary:        {{DESC}}
Name:           {{SHORT_NAME}}
Version:        {{VERSION}}
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/{{SHORT_NAME}}

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:      checksum.sha512

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.19

Requires:       systemd

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
{{DESC}}.

################################################################################

%prep
%{crc_check}

%setup -q

%build
if [[ ! -d "%{name}/vendor" ]] ; then
  echo "This package requires vendored dependencies"
  exit 1
fi

pushd %{name}
  %{__make} %{?_smp_mflags} all
popd

%install
rm -rf %{buildroot}

install -dDm 755 %{buildroot}%{_bindir}
install -dDm 755 %{buildroot}%{_sysconfdir}/logrotate.d
install -dDm 755 %{buildroot}%{_unitdir}
install -dDm 755 %{buildroot}%{_logdir}/%{name}

install -pm 755 %{name}/%{name} \
                %{buildroot}%{_bindir}/

install -pm 644 %{name}/common/%{name}.knf \
                %{buildroot}%{_sysconfdir}/

install -pm 644 %{name}/common/%{name}.logrotate \
                %{buildroot}%{_sysconfdir}/logrotate.d/%{name}

install -pm 644 %{name}/common/%{name}.service \
                %{buildroot}%{_unitdir}/

%clean
rm -rf %{buildroot}

%pre
getent group %{name} >/dev/null || groupadd -r %{name}
getent passwd %{name} >/dev/null || useradd -r -M -g %{name} -s /sbin/nologin %{name}
exit 0

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%attr(-,%{name},%{name}) %dir %{_logdir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}.knf
%config(noreplace) %{_sysconfdir}/logrotate.d/%{name}
%{_unitdir}/%{name}.service
%{_bindir}/%{name}

################################################################################

%changelog
* {{SPEC_CHANGELOG_DATE}} Anton Novojilov <andy@essentialkaos.com> - {{VERSION}}-0
- Initial build for kaos-repo
