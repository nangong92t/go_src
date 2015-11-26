#!/usr/bin/perl -w -T

use strict;
use RPC::PlClient;
use Data::Dumper;

$ENV{PATH} = '/bin:/usr/bin';
delete @ENV{qw(IFS CDPATH ENV BASH_ENV)};
delete @ENV{qw(HOME LOGNAME MAIL USER USERNAME)};

package SyncServer;
use IPC::Run qw(run);
use vars qw($VERSION);
$VERSION = '2.0';

# the user and group under which the regular operations
# (sources checkout, sources export and configuration, interogation)
# are made
use constant SYNC_USER             => 'git-sync';
use constant SYNC_GROUP            => 'git-sync';
use constant SYNC_GROUP_ADDITIONAL => 'git-sync';
use constant GET_SYNC_USER_HOME    => sub {
    (getpwnam(SYNC_USER))[7];
};

# git binaries
use constant GITVERSION_PATH  => '/usr/bin/git';
use constant GITCLIENT_PATH   => '/usr/bin/git';

my %CONF = (
    web_devel => {
        checked_out_dir      => 'web_devel',
        exported_dir         => 'web_devel',
        branch               => 'devel',
        restartable_services => [],
        reloadable_services  => [],
    },
    web_product => {
        checked_out_dir      => 'web',
        exported_dir         => 'web',
        branch               => 'master',
        restartable_services => [{
            name      => 'Memcached',
            ctl_path  => '/usr/local/apache/bin/apachectl',
            pid_file  => '/usr/local/apache/logs/httpd.pid',
            wait_time => 30, # seconds
        }],
        reloadable_services  => [],
    },
    web_publish => {
        checked_out_dir      => 'web_publish',
        exported_dir         => 'web_publish',
        branch               => 'publish',
        restartable_services => [],
        reloadable_services  => [],
    },
    static_product => {
        checked_out_dir      => 'static',
        exported_dir         => 'static',
        branch               => 'master',
        restartable_services => [],
        reloadable_services  => [],
    },
    static_devel => {
        checked_out_dir      => 'static',
        exported_dir         => 'static',
        branch               => 'devel',
        restartable_services => [],
        reloadable_services  => [],
    },
    service_product => {
        checked_out_dir      => 'service',
        exported_dir         => 'service',
        branch               => 'master',
        restartable_services => [],
        reloadable_services  => [],
    },
    service_devel => {
        checked_out_dir      => 'service_devel',
        exported_dir         => 'service_devel',
        branch               => 'devel',
        restartable_services => [],
        reloadable_services  => [],
    },
    vendor_product => {
        checked_out_dir      => 'vendor',
        exported_dir         => 'vendor',
        branch               => 'master',
        restartable_services => [{
            name      => 'VendorService',
            is_root   => 1,
            ctl_path  => '/usr/local/git-sync/vendor/bin/workermand',
            pid_file  => '/usr/local/git-sync/vendor/conf/bin/pid',
            wait_time => 30, # seconds
        }],
        reloadable_services  => [],
    },
    vendor_devel => {
        checked_out_dir      => 'vendor_devel',
        exported_dir         => 'vendor_devel',
        branch               => 'devel',
        restartable_services => [{
            name      => 'VendorService',
            is_root   => 1,
            ctl_path  => '/usr/local/git-sync/vendor/bin/workermand',
            pid_file  => '/usr/local/git-sync/vendor/conf/bin/pid',
            wait_time => 30, # seconds
        }],
        reloadable_services  => [],
    },
    maintance => {
        checked_out_dir      => 'maintance',
        exported_dir         => 'maintance',
        branch               => 'master',
        restartable_services => [{
            name      => 'StaticServer',
            is_root   => 1,
            ctl_path  => '/usr/local/git-sync/maintance/staticserver/static-server',
            pid_file  => '/usr/local/git-sync/maintance/sync.pid',
            wait_time => 30, # seconds
        }],
        reloadable_services  => [],
    },

);
$CONF{webmail_servers} = $CONF{web_servers};


my @saved_uid = ($<, $>);
my @saved_gid = ($(, $));

sub new {
    bless { 'server' => $RPC_SyncServer::server }, shift;
}

sub lose_privileges {
    my $uid = getpwnam(SYNC_USER);
    my $gid = getgrnam(SYNC_GROUP);
    my $gid_additional = join(" ", map {
            scalar getgrnam($_);
        } split(/\s+/, SYNC_GROUP_ADDITIONAL));
    $gid_additional = $gid if $gid_additional =~ /^\s*$/;
    $( = $gid;
    $) = "$gid $gid_additional";
    $> = $< = $uid;
}

sub restore_privileges {
    ($<, $>) = @saved_uid;
    ($(, $)) = @saved_gid;
}

sub lose_root_privileges {
    my $uid = getpwnam('root');
    my $gid = getgrnam('root');
    my $gid_additional = join(" ", map {
            scalar getgrnam($_);
        } split(/\s+/, SYNC_GROUP_ADDITIONAL));
    $gid_additional = $gid if $gid_additional =~ /^\s*$/;
    $( = $gid;
    $) = "$gid $gid_additional";
    $> = $< = $uid;
}

sub get_head_revision {
    my ($self, $conf_type) = @_;
    $self->{server}->Debug('get_head_revision');

    my $co_dir = $CONF{$conf_type}->{checked_out_dir};
    return (undef, "No 'checked_out_dir' defined for conf_type '$conf_type'")
    unless $co_dir;

    lose_privileges();

    # svn status -u ~git-sync/co/empty-for-svnversion
    my ($out, $err);
    local $SIG{PIPE} = 'IGNORE';
    run [ GITCLIENT_PATH, 'status', '-u',
    join('/', GET_SYNC_USER_HOME->(), $co_dir, 'empty-for-svnversion') ],
    \undef, \$out, \$err;

    restore_privileges();

    if ($out =~ /^Status against revision:\s*(\d+)\s*$/) {
        $out = $1;
    }
    return ($out, $err);
}

sub get_production_revision {
    my ($self, $conf_type) = @_;
    $self->{server}->Debug('get_production_revision');

    my $exported_dir = $CONF{$conf_type}->{exported_dir};
    return (undef, "No 'exported_dir' defined for conf_type '$conf_type'")
    unless $exported_dir;

    lose_privileges();

    # ~git-sync/export/production is symlink to ~git-sync/export/r-revision
    my $revision;
    my $path = join('/', GET_SYNC_USER_HOME->(), $exported_dir, '');
    if (-l $path) {
        $revision = readlink $path;
        $revision =~ s/^r-(\d+)$/$1/;
        undef $revision if $revision !~ /^\d+$/;
    }

    restore_privileges();

    return $revision;
}

sub get_last_changed_revision {
    my ($self, $conf_type) = @_;
    $self->{server}->Debug('get_last_changed_revision');

    my $co_dir = $CONF{$conf_type}->{checked_out_dir};
    return (undef, "No 'checked_out_dir' defined for conf_type '$conf_type'")
    unless $co_dir;

    lose_privileges();

    # svnversion -c ~git-sync/co
    my ($out, $err, $revision);
    local $SIG{PIPE} = 'IGNORE';
    run [ GITVERSION_PATH, '-cn', join('/', GET_SYNC_USER_HOME->(), $co_dir) ],
    \undef, \$out, \$err;
    ($revision) = ($out =~ /^\d+:(\d+)/);

    restore_privileges();

    return ($revision, $err);
}

sub get_checked_out_revision {
    my ($self, $conf_type) = @_;
    $self->{server}->Debug('get_checked_out_revision');

    my $co_dir = $CONF{$conf_type}->{checked_out_dir};
    return (undef, "No 'checked_out_dir' defined for conf_type '$conf_type'")
    unless $co_dir;

    lose_privileges();

    # git --git-dir ~git-sync/co/.git log --pretty=format:"%h - %an, %ar : %s" -1
    my ($out, $err, $revision);
    local $SIG{PIPE} = 'IGNORE';
    run [ GITVERSION_PATH, '--git-dir', join('/', GET_SYNC_USER_HOME->(),
        $co_dir, '.git'), 'log', '--pretty=format:"%h - %an, %ar : %s"', '-1'],
    \undef, \$out, \$err;

    $revision = $out;

    restore_privileges();

    return ($revision, $err);
}

sub check_out_sources {
    my ($self, $revision, $conf_type) = @_;
    $self->{server}->Debug('check_out_sources');

    my $co_dir = $CONF{$conf_type}->{checked_out_dir};
    my $branch = $CONF{$conf_type}->{branch};

    return (undef, "No 'checked_out_dir' defined for conf_type '$conf_type'")
    unless $co_dir;

    lose_privileges();
    my $orig_umask = umask 0002;

    # git --git-dir ~git-sync/co/.git pull origin master
    # pushd /usr/local/git-sync/web && git pull origin master && popd
    # sh -c 'cd /usr/local/git-sync/co && git pull origin master'
    my @cmd = ('/bin/sh', '-c', 'cd '.join('/', GET_SYNC_USER_HOME->(), $co_dir) . ' && ' . GITCLIENT_PATH . ' pull origin ' . $branch.':'.$branch);
    @cmd = ('/bin/sh', '-c', 'cd '.join('/', GET_SYNC_USER_HOME->(), $co_dir) . ' && ' . GITCLIENT_PATH .' reset --hard '. $1)
    if defined $revision and $revision =~ /^\s*(\w+)\s*$/;
    my ($out, $err);
    local $SIG{PIPE} = 'IGNORE';
    run \@cmd, \undef, \$out, \$err;

    umask $orig_umask;
    restore_privileges();

    return ($out, $err);
}

sub export_sources {
    my ($self, $conf_type) = @_;
    $self->{server}->Debug('export_sources');

    my $co_dir = $CONF{$conf_type}->{checked_out_dir};
    return (undef, "No 'checked_out_dir' defined for conf_type '$conf_type'")
    unless $co_dir;

    my $exported_dir = $CONF{$conf_type}->{exported_dir};
    return (undef, "No 'exported_dir' defined for conf_type '$conf_type'")
    unless $exported_dir;

    lose_privileges();
    my $orig_umask = umask 0002;

    # ~git-sync/co/scripts/export.sh
    my ($out, $err);
    local $SIG{PIPE} = 'IGNORE';
    run [ join('/', GET_SYNC_USER_HOME->(), $co_dir, 'scripts', 'export.sh'),
    $co_dir, $exported_dir ], \undef, \$out, \$err;

    umask $orig_umask;
    restore_privileges();

    return ($out, $err);
}

sub restart_services {
    my ($self, $conf_type) = @_;

    return (undef, "No 'restartable_services' defined for conf_type '$conf_type'")
    unless exists $CONF{$conf_type}->{restartable_services};

    my ($all_out_stop, $all_err_stop, $all_out_start, $all_err_start);

    foreach my $service (@{ $CONF{$conf_type}->{restartable_services} }) {
        local $SIG{PIPE} = 'IGNORE';

        lose_root_privileges() if $service->{is_root};
        my ($out_stop, $err_stop, $out_start, $err_start);
        run [ $service->{ctl_path}, 'stop' ], \undef, \$out_stop, \$err_stop;
        $self->{server}->Log(notice => "Stopped $service->{name}");

        eval {
            local $SIG{ALRM} = sub { die "alarm" };
            alarm $service->{wait_time};
            while (-e $service->{pid_file}) {
                select(undef, undef, undef, 0.1);
            }
            alarm 0;
        };
        if ($@ and $@ !~ /^alarm/) {
            $self->{server}->Log(crit => "Error stopping $service->{name}: $@");
        } else {
            local $SIG{PIPE} = 'IGNORE';
            if (run([ $service->{ctl_path}, 'start' ], \undef, \$out_start, \$err_start)) {
                $self->{server}->Log(notice => "Started $service->{name}");
            } else {
                $self->{server}->Log(notice => "Failed to start $service->{name}: $?");
            }
        }

        $all_out_stop .= "$service->{name}: $out_stop\n" if length $out_stop;
        $all_err_stop .= "$service->{name}: $err_stop\n" if length $err_stop;
        $all_out_start .= "$service->{name}: $out_start\n" if length $out_start;
        $all_err_start .= "$service->{name}: $err_start\n" if length $err_start;

        restore_privileges() if $service->{is_root};

    }

    return ($all_out_stop, $all_err_stop, $all_out_start, $all_err_start);
}

sub reload_services {
    my ($self, $conf_type) = @_;

    return (undef, "No 'reloadable_services' defined for conf_type '$conf_type'")
    unless exists $CONF{$conf_type}->{reloadable_services};

    my ($all_out_reload, $all_err_reload);

    foreach my $service (@{ $CONF{$conf_type}->{reloadable_services} }) {
        local $SIG{PIPE} = 'IGNORE';

        my ($out_reload, $err_reload);
        run [ $service->{ctl_path}, 'reload' ], \undef, \$out_reload, \$err_reload;
        $self->{server}->Log(notice => "Reloaded $service->{name}");

        $all_out_reload .= "$service->{name}: $out_reload\n" if length $out_reload;
        $all_err_reload .= "$service->{name}: $err_reload\n" if length $err_reload;
    }

    return ($all_out_reload, $all_err_reload);
}

sub do_export {
    my ($self, $params) = @_;
    my (@output, $out, $err);

    unless (defined $params and exists $params->{conf_type}) {
        return ({
                method => 'do_export',
                out => '',
                err => "Missing 'conf_type' parameter",
            });
    }

    my $revision = (exists($params->{revision}) and $params->{revision} =~ /^\w+$/) ? $params->{revision} : undef;
    if (!$params->{build_static_page}) {
        ($out, $err) = $self->check_out_sources($revision, $params->{conf_type});
        push @output, {
            method => 'check_out_sources',
            out => $out,
            err => $err,
        };
    }

    my $isdev = exists($params->{isdev}) ? $params->{isdev} : undef;

    unless ($isdev) {
        ($out, $err) = $self->export_sources($params->{conf_type});
        push @output, {
            method => 'export_sources',
            out => $out,
            err => $err,
        };
    }

    if (exists $params->{restart_services} and $params->{restart_services}) {

        my ($out_stop, $err_stop,
            $out_start, $err_start) = $self->restart_services($params->{conf_type});
        push @output, {
            method => 'stop_services',
            out => $out_stop,
            err => $err_stop,
        };
        push @output, {
            method => 'start_services',
            out => $out_start,
            err => $err_start,
        };
    }

    if (exists $params->{reload_services} and $params->{reload_services}) {
        my ($out_reload, $err_reload) = $self->reload_services($params->{conf_type});
        push @output, {
            method => 'reload_services',
            out => $out_reload,
            err => $err_reload,
        };
    }

    if (exists $params->{clean_page_template} and $params->{clean_page_template} == 1) {
        my $co_dir = $CONF{$params->{conf_type}}->{checked_out_dir};
        my $curDir = GET_SYNC_USER_HOME->() . '/' . $co_dir . '/apps';
        my @apps   = ('wechat', 'api', 'pc', 'admin');

        foreach my $curF (@apps) { 
            run "rm -rf $curDir/$curF/protected/runtime/*";
        }
        
        push @output, {
            method => 'clean_page_template',
            out => 'cleaned page template done!',
        };
    }

    # update web static file version for CDN.
    if (exists $params->{update_static_version_forcdn} and $params->{update_static_version_forcdn} == 1) {
        my $out = {
            method => 'update_static_version_forcdn ',
            out => 'update web static url version done!',
        };
        my $version= (exists($params->{version}) and $params->{version} =~ /^\w+$/) ? $params->{version} : undef;
        if (!$version) {
            $out->{'out'}   = 'update web static url version failed: no find current static version!';
            push @output, $out;
        }

        my $co_dir = $CONF{$params->{conf_type}}->{checked_out_dir};
        my $curDir = GET_SYNC_USER_HOME->() . '/' . $co_dir . '/config';

        my $file    = "$curDir/main.php";
        if (-e "$curDir/development.php") {
            $file   = "$curDir/development.php";
        }

    	lose_privileges();

        open Config, "<$file";
        my $content = join "", <Config>;
        close Config;
        #run [ '/bin/sed', '-i', 's/"StaticVersion".*=>.*".*"/"StaticVersion" => "'.$version.'"/', $file ], \undef, \$out_reload, \$err_reload;
        $content    =~ s#'StaticVersion'.*=>.*'.*'#'StaticVersion' => '$version'#;
        open Config, ">$file" or $out->{'out'} = 'update web static url version failed: '. $!;
        print Config $content;
        close Config;

    	restore_privileges();

        push @output, $out;
    }

    if (exists $params->{build_static_page} and $params->{build_static_page} == 1) {
        my $co_dir = $CONF{$params->{conf_type}}->{checked_out_dir};
        my $curDir = GET_SYNC_USER_HOME->() . '/' . $co_dir . '/apps';

        run "/usr/bin/wget http://show.wepiao.com/index.php -O $curDir/index.html";
        
        push @output, {
            method => 'build_static_page',
            out => 'build index.html done!',
        }
    }

    return @output;
}

sub do_update_vost_dns{
    my ($self, $params) = @_;
    $self->{server}->Log(notice => "Updating vost DNS setting...");
    my $msg = run "perl /usr/local/git-sync/co/scripts/cron/vost/update_dns.pl" ;
    $self->{server}->Log(notice => "Updated vost DNS setting successful.") if ($msg=~/OK/ig);
    return [$msg];
}

sub get_csb_context {
    my $rpc_client = RPC::PlClient->new(peeraddr => '127.0.0.1', peerport => 2626, application => 'RPC_CSB_Server', version => '1.0', logfile => 1);
    my $rpc_context = $rpc_client->ClientObject('CSB_Server', 'new');
    return $rpc_context;
}

sub do_setup_clone_site {
    my ($self, $params) = @_;

    $self->{server}->Log(notice => 'Setting up clone site ...');

    my $rpc_context = &get_csb_context;
    my ($out, $err) = $rpc_context->do_setup_clone_site($params);
    return {
        'method' => 'do_setup_clone_site',
        'out' => $out,
        'err' => $err,
    };
}

package RPC_SyncServer;
require RPC::PlServer;
use IO::File;
use POSIX qw(setsid);
use vars qw($VERSION @ISA $server);
$VERSION = '1.0';
@ISA = qw(RPC::PlServer);

# rpc server
use constant PID_FILE        => '/var/run/syncserver.pid';
use constant LOG_FILE        => '/var/log/syncserver';
use constant PORT            => 26;

sub daemonize {
    chdir '/'               or die "Can't chdir to /: $!";
    open STDIN, '/dev/null' or die "Can't read /dev/null: $!";
    open STDOUT, '>/dev/null'
        or die "Can't write to /dev/null: $!";
    defined(my $pid = fork) or die "Can't fork: $!";
    exit if $pid;
    setsid                  or die "Can't start a new session: $!";
    open STDERR, '>&STDOUT' or die "Can't dup stdout: $!";
}

eval {
    daemonize();
    my $logfile = IO::File->new(LOG_FILE, 'a')
        or die "Couldn't set logfile: $!";
    $logfile->autoflush(1);

    $server = RPC_SyncServer->new({
            'pidfile'    => PID_FILE,
            'logfile'    => $logfile,
            'debug'      => 0,
            'user'       => 0, # root, but due to a bug in Net::Daemon 0.37 the
            # server dies with an error when the name is used
            'group'      => 0, # root
            'localport'  => PORT,
            'mode'       => 'single', # only one connection at a time
            'clients'    => [
            {
                'mask'   => '^127\.0\.0\.1$',
                'accept' => 1,
            },
            {
                'mask'   => '^10\.251\.233\.145$',
                'accept' => 1,
            },
            {
                'mask'   => '^182\.254\.232\.149$',
                'accept' => 1,
            },
            # Deny everything else
            {
                'mask'   => '.*',
                'accept' => 0,
            },
            ],
            'methods'    => {
                'RPC_SyncServer' => {
                    'ClientObject' => 1,
                    'CallMethod'   => 1,
                    'NewHandle'    => 1,
                },
                'SyncServer' => {
                    'new'                       => 1,
                    'get_head_revision'         => 1,
                    'get_production_revision'   => 1,
                    'get_checked_out_revision'  => 1,
                    'get_last_changed_revision' => 1,
                    'check_out_sources'         => 1,
                    'export_sources'            => 1,
                    'restart_services'          => 1,
                    'reload_services'           => 1,
                    'do_export'                 => 1,
                    'do_update_vost_dns'        => 1,
                    'do_setup_clone_site'       => 1,
                },
            },
        });
    $SIG{PIPE} = 'IGNORE';
    $server->Bind();
};
if ($@) {
    print STDERR "Couldn't create RPC_SyncServer instance: $@"; # 'emacs
}
