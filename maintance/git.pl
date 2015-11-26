#!/usr/bin/perl -w
# $Id$

use strict;
use Data::Dumper;
use CGI;
use MIME::Base64;
require RPC::PlClient;
use Parallel::Jobs qw(start_job watch_jobs);
use Storable ();
use Socket;


my %servers = (
    'Web Product' => {
        conf_type               => 'web_product',
        restartable_services => 0,
        servers => {
           'TX02' => { ip => '10.251.58.96', isdev => 1 },
           'TX11' => { ip => '10.251.46.47', isdev => 1 },
           'TX12' => { ip => '10.251.89.164', isdev => 1 },
        },
    },
    'Web Publish' => {
        conf_type               => 'web_publish',
        restartable_services => 0,
        servers => {
           'TX11' => { ip => '10.251.46.47', isdev => 1 },
           'TX12' => { ip => '10.251.89.164', isdev => 1 },
        },
    },
    'Web Devel' => {
        conf_type               => 'web_devel',
        restartable_services => 0,
        servers => {
           'Task' => { ip => '111.203.241.199', isdev => 1 },
        },
    },
    'Service Product' => {
        conf_type               => 'service_product',
        restartable_services => 0,
        servers => {
            'TX13' => { ip => '10.104.4.66', isdev => 1 },
            'TX14' => { ip => '10.250.136.54', isdev => 1 },
        },
    },
    'Service Devel' => {
        conf_type               => 'service_devel',
        restartable_services => 0,
        servers => {
           'Task' => { ip => '111.203.241.199', isdev => 1 },
        },
    },
    'Static Product' => {
        conf_type               => 'static_product',
        restartable_services => 0,
        servers => {
            'TX17' => { ip => '10.251.59.45', isdev => 1 },
            'TX18' => { ip => '10.251.58.71', isdev => 1 },
        },
        update_version_group => 'Web Product'
    },
    'Static Devel' => {
        conf_type               => 'static_devel',
        restartable_services => 0,
        servers => {
            'Task' => { ip => '111.203.241.199', isdev => 1 },
        },
        update_version_group => 'Web Devel'
    },
    'Vendor Product' => {
        conf_type               => 'vendor_product',
        restartable_services => 1,
        servers => {
           'TX13' => { ip => '10.104.4.66', isdev => 1 },
           'TX14' => { ip => '10.250.136.54', isdev => 1 },
        },
    },
    'Vendor Devel' => {
        conf_type               => 'vendor_devel',
        restartable_services => 1,
        servers => {
           'Task' => { ip => '111.203.241.199', isdev => 1 },
        },
    },
    'Synchronizing' => {
        conf_type               => 'maintance',
        restartable_services => 1,
        servers => {
            'Task' => { ip => '111.203.241.199', isdev => 1 },
            'TX11' => { ip => '10.251.46.47', isdev => 1 },
            'TX12' => { ip => '10.251.89.164', isdev => 1 },
            'TX13' => { ip => '10.104.4.66', isdev => 1 },
            'TX14' => { ip => '10.250.136.54', isdev => 1 },
            'TX17' => { ip => '10.251.59.45', isdev => 1 },
            'TX18' => { ip => '10.251.58.71', isdev => 1 },
        },
    },
);

use constant DEFAULT_STATIC_SERVER  => 'tx15';
use constant DEFAULT_RPC_SERVER     => '182.254.232.149';

use constant DEFAULT_CSS => '* { font-size:14px; }
a { color:#009; text-decoration:none; }
a:hover { color:#eee; background:#369; }
label { display:block; font-weight:bold; }
.matrix label { background-color:#eee; cursor:pointer; display:block; float:left; font-weight:normal; margin:0.2em; padding:0.2em 0; width:19%; }
.matrix label:hover { background-color:#369; color:#eee; }
.clear { clear:both; }
.form { margin:1em 0; }
.form input { border:1px solid #ccc; padding:.2em; font-size:16px; }
.error { color:red; }
.obvious { font-weight:bold; }';


my %status;
$| = 1;

my $q = new CGI;
my $mode = $q->param('rm');

print $q->header;
if (defined $mode and $mode eq 'do_sync') {
    do_sync_page($q);
}
else {
    status_page($q);
}

sub do_sync_server {
    my ($server, $params) = @_;

    my @output;
    eval {
        local $SIG{ALRM} = sub { die "timed out" };
        alarm 300;
        my $rpc_context = get_rpc_context($params->{conf}->{servers}->{$server}->{ip});

        $params->{isdev}     = $params->{conf}->{servers}->{$server}->{isdev};
        $params->{conf_type} = $params->{conf}->{conf_type};

        @output = $rpc_context->do_export($params);
        alarm 0;
    };
    return [], $@ if $@;
    return \@output;
}

sub start_do_sync_server {
    my ($server, $params) = @_;

    start_our_job(
                  sub {
                      my $pid = shift;
                      $status{$pid} = { server => $server };
                  },
                  sub {
                      my ($output, $error) = do_sync_server($server, $params);
                      print STDOUT Storable::freeze($output) if @$output;
                      print STDERR $error if defined $error;
                      exit 0;
                  }
                  );
}

sub display_sync_output {
    my ($q, $info) = @_;

    my $server = $info->{server};
    print $q->b("$server:"), $q->br, "\n";
    if (exists $info->{STDERR} and length $info->{STDERR}) {
        print $q->div({ -style => 'color: red' }, $info->{STDERR}), "\n";
    } else {
        my @output = @{ Storable::thaw($info->{STDOUT}) };
        foreach (@output) {
            if (ref $_) {
                print $_->{method}, "\n";
                print $q->pre($_->{out}), "\n";
                print $q->pre({ -style => 'color: red' }, $_->{err}), "\n"
                    if defined $_->{err} and length $_->{err};
            } else {
                print $_, "\n";
            }
        }
    }
    print $q->hr, "\n";
}

sub validate_clone_site {
    my ($matrix, $domain) = @_;

    if (length $matrix == 0) {
        return 'Matrix name is required!';
    }
    elsif (length $domain == 0) {
        return 'Domain name is required!';
    }
    elsif ($matrix !~ m/^[-_\w]+\.\w{2,4}(\.\w{2})?$/) {
        return 'Matrix name is invalid!';
    }
    elsif ($domain !~ m/^[-_\w]+\.\w{2,4}(\.\w{2})?$/) {
        return 'Domain name is invalid!';
    }
    elsif (lc($matrix) eq lc($domain)) {
        return 'Are you kidding?';
    }

}

sub get_rpc_context {
    my ($peeraddr) = @_;

    my $rpc_client = RPC::PlClient->new(peeraddr => $peeraddr, peerport => 26, application => 'RPC_SyncServer', version => '1.0', logfile => 1);
    my $rpc_context = $rpc_client->ClientObject('SyncServer', 'new');
    return $rpc_context;
}

sub do_sync_page {
    my $q = shift;

    print $q->start_html('Synchronizing servers page'), "\n";
    print $q->h3('Synchronizing servers'), "\n";

    my $group = $q->param('group');
    unless (defined $group) {
        print $q->div({ -style => 'color: red' }, "The group of servers to synchronize is missing. Please check the URL you are using."), "\n";
        print $q->end_html, "\n";
        return;
    }

    my $end_page = $q->br . $q->br .
        $q->a({-href => $q->url . '?group=' . $group }, "Back to the status page of '$group'") . $q->br .
        $q->a({-href => $q->url }, 'Back to the server type selection page') . $q->br .
        "\n" . $q->end_html . "\n";

    my $conf = $servers{$group};

    my @servers_to_synchronize = ();
    foreach my $param_name ($q->param) {
        if ($param_name =~ /^server_([\w+.-]+)$/i && $1 ne 'all') {
            my $server = $1;
            if (defined $q->param("server_$server") and
                $q->param("server_$server") == 1) {
                push @servers_to_synchronize, $server;
            }
        }
    }

    @servers_to_synchronize = sort sort_servers @servers_to_synchronize;
    my @servers_order = map { { server => $_ } } @servers_to_synchronize;

    if (@servers_to_synchronize == 0) {
        print $q->div({ -style => 'color: red' }, "No server has been selected. Please select one or more servers to synchronize and try again."), "\n";
        print $end_page;
        return;
    }

    my $revision;
    unless (defined $q->param('revision_head') and $q->param('revision_head') == 1) {
        if (defined $q->param('revision') and
            $q->param('revision') =~ /^(\w+)$/) {
            $revision = $1;
        } else {
            print $q->div({ -style => 'color: red' }, "The revision number is missing or is not a number. Please enter a revision number into the specific revision field or check the \"Repository head revision\" box and try again."), "\n";
            print $end_page;
            return;
        }
    }

    my $restart_services = (defined $q->param('restart_services') and
                            $q->param('restart_services') == 1);
    my $reload_services = (defined $q->param('reload_services') and
                            $q->param('reload_services') == 1);

    my $clean_page_template = (defined $q->param('clean_page_template') and
			    $q->param('clean_page_template') == 1);

    my $build_static_page = (defined $q->param('build_static_page') and
			    $q->param('build_static_page') == 1);

    my $fast_sync = (defined $q->param('fast_sync') and
                     $q->param('fast_sync') == 1);

    foreach my $server (@servers_to_synchronize) {
        start_do_sync_server($server,
                             {  
                                revision => $revision,
                                restart_services => $restart_services,
                                reload_services => $reload_services,
                                conf             => $conf,
                                fast_sync => $fast_sync,
 			        clean_page_template => $clean_page_template,
			        build_static_page  => $build_static_page,
                             },
        );
    }

    my $info;
    watch_our_jobs(\@servers_order, ['server'], sub {
        $info = shift;

        display_sync_output($q, $info);
    });

    # to prepare the update cdn version service
    if (defined $q->param('update_static_group')) {
        update_static_version_forcdn($q->param('update_static_group'), $info);
    }

    print $end_page;
}

sub update_static_version_forcdn {
    my ($groupName, $info) = @_;
    sleep(3);

    my $version = $1 if (defined $info->{STDOUT} and $info->{STDOUT} =~ m/\w{7}..(\w{7}) /);
    return 0 if !$version;
    print "<p><b>Updating Web Static Version For CDN.</b></p>";
    my @servers_order;
    my $curServers   = $servers{ $groupName };
    return 0 if !$curServers;
    foreach my $server (keys %{ $curServers->{'servers'} } ) {
        start_do_sync_server($server,
                             {  
                                revision => undef,
                                version => $version,
                                conf             => $curServers,
                                fast_sync => 0,
                                clean_page_template => 1,
			        update_static_version_forcdn	=> 1
                             }
        );
        push @servers_order, { server => $server };
    }

    watch_our_jobs(\@servers_order, ['server'], sub {
        my $info2 = shift;
        $info2->{STDOUT} =~ s/12345678//g;
        $info2->{STDOUT} =~ s/\n/<br\/>/g;
        print "<b>".$info2->{server} . ":</b><br/>";
        print '<pre>' . $info2->{STDOUT} .'</pre>';
    });
}


sub get_revision {
    my ($server, $conf, $revision_type) = @_;
    my ($revision, $error);
    eval {
        local $SIG{ALRM} = sub { die "timed out" };
        alarm 15;
        my $rpc_context = get_rpc_context($conf->{servers}->{$server}->{ip});
        my $method = "get_${revision_type}_revision";
        ($revision, $error) = $rpc_context->$method($conf->{conf_type});
        alarm 0;
    };
    $error .= $@ if $@;
    return ($revision, $error);
}

sub start_get_revision {
    my ($server, $conf, $revision_type) = @_;
    start_our_job(
                  sub {
                      my $pid = shift;
                      $status{$pid} = { server => $server, revision_type => $revision_type };
                  },
                  sub {
                      my ($revision, $error) = get_revision($server, $conf, $revision_type);
                      print STDOUT $revision if defined $revision;
                      print STDERR $error if defined $error;
                      exit 0;
                  }
                  );
}

sub revision_as_text {
    my ($revision, $error) = @_;
    my $output;
    if (!defined $revision && !defined $error) {
        $output = $q->span({ -style => 'color: red' }, 'unknown');
    } else {
        $output = $revision;
        $output .= $q->span({ -style => 'color: red' }, $error)
            if defined $error and length $error > 0;
    }
    return $output;
}

sub display_server_status {
    my ($q, $group, $info) = @_;

    my $server_info = $servers{$group}->{servers}->{ $info->{server} };
    my $is_default_server = !exists $server_info->{default} || $server_info->{default};
    my $css_class = $is_default_server ? 'def_server' : 'ndef_server';
    my $checked = scalar( keys %{ $servers{$group}->{servers} } ) <= 1;

    print $q->Tr($q->td([ ${info}->{server},
                          $q->div({ -align => 'center' }, revision_as_text($info->{STDOUT}, $info->{STDERR})),
                          $q->div({ -align => 'center', -class => $css_class },
                                  $q->checkbox(-name => "server_$info->{server}",
                                               -checked => $checked,
                                               -value => 1,
                                               -label => 'yes',
                                               -onClick => "javascript: check_sync_server(this)")
                          ),
                        ])),
    "\n";
}

sub status_page {
    my $q = shift;

    my $css = <<'END';
* { font-size:14px; }
a { color:#009; text-decoration:none; }
a:hover { background:#ebdde2; }
.options ul { list-style-type:none; margin:0; padding:0; }
.options li { margin:.2em; float:left; width:19%; }
.options a { display:block; padding:.5em; color:#009; background:#eee; text-decoration:none; }
.options a:hover { color:#eee; background:#369; }
.options .clone { font-weight:bold; }
END

    my $group = $q->param('group');
    unless (defined $group) {
        print $q->start_html(-title => 'Status page', -style => { -code => $css }, );
        print $q->h2('Select the type of servers to synchronize:');
        print $q->start_div({-class => 'options'});
        print $q->start_ul;
        print $q->li([ map { $q->a({-href => $q->url . "?group=$_" }, $_) } sort keys %servers ]);
        print $q->end_ul;
        print $q->end_div;
        print $q->end_html;
        return;
    }

    my $conf = $servers{$group};

    my $jscript = <<'END';
// uncheck 'All default servers' checkbox when a specific server is checked
function check_sync_server(server) {
    if (server.checked == 0) {
        document.forms.synchronize.server_all.checked = 0;
    }
}

// if 'All default servers' is checked, then
// check each default server
// and uncheck each non-default server
function check_all_servers(server_all) {
    if (server_all.checked) {
        var i;
        for (i = 0; i < document.forms.synchronize.elements.length; i++) {
            var e = document.forms.synchronize.elements[i];
            if (e.name.match(/^server_/) && e.name != "server_all") {
                var server = e.name.substr(7);
                var is_default_server = true;
                for (var s in nd_servers) {
                    if (s == server) {
                        is_default_server = false;
                        break;
                    }
                }
                e.checked = is_default_server;
            }
        }
    } else {
        var j;
        for (j = 0; j < document.forms.synchronize.elements.length; j++) {
            var e = document.forms.synchronize.elements[j];
            if (e.name.match(/^server_/) && e.name != "server_all") {
                e.checked = false;
            }
        }
    }
}

// servers to not sync by default
END

    $jscript .= "var nd_servers = new Array(";
    $jscript .= join(", ", map { '"' . $_ . '"' } grep { exists $conf->{servers}->{ $_ }->{default} && !$conf->{servers}->{ $_ }->{default} } keys %{ $conf->{servers} });
    $jscript .= ");\n";

    $css = <<'END';
<!--
div.ndef_server {
  color: gray;
}

div.def_server {
}
-->
END

    print $q->start_html(-title  => 'Status page',
                         -script => $jscript,
                         -style  => { -code => $css }, ), "\n";

    print $q->h3('Servers status'), "\n";

    my @servers_order = map {
        { server => $_,
          revision_type => defined($conf->{servers}->{$_}->{isdev}) ? 'checked_out' : 'production' }
    } sort sort_servers keys %{ $conf->{servers} };

    foreach (@servers_order) {
        start_get_revision($_->{server}, $conf, $_->{revision_type});
    }

    my $last_changed_rev_server = $conf->{last_changed_rev_server};
    if (defined $last_changed_rev_server) {
        start_get_revision($last_changed_rev_server, $conf, 'last_changed');
    }

    print $q->start_form(-method => 'POST', -name => 'synchronize'), "\n";
    print $q->hidden(rm => 'do_sync'), "\n";
    print $q->hidden(group => $group), "\n";

    print $q->start_table({ -border => 1 }), "\n";
    print $q->Tr($q->th([ 'Server', 'Production revision', 'Synchronize?' ])), "\n";

    print $q->Tr($q->td([ '', '',
                          $q->checkbox(-name => 'server_all',
                                       # if only one server, then check this
                                       -checked => @servers_order <= 1,
                                       -value => 1, -label => 'All default servers',
                                       -onClick => 'check_all_servers(this)')
                        ])), "\n";


    watch_our_jobs(\@servers_order, [qw(server revision_type)], sub {
        my $info = shift;
        display_server_status($q, $group, $info);
    });

    print $q->end_table, $q->br, "\n";

    if (defined $last_changed_rev_server) {
        print $q->b('Last changed revision of trunk: ');
        foreach my $pid (keys %status) {
            my $info = $status{$pid};
            if ($info->{server} eq $last_changed_rev_server and $info->{revision_type} = 'last_changed') {
                print revision_as_text($info->{STDOUT}, $info->{STDERR});
                last;
            }
        }
    }

    # to which revision
    print $q->b('To which revision?'), $q->br, "\n";
    print $q->checkbox(-name => 'revision_head', -checked => 1,
               -value => 1, -label => 'Repository head revision',
               -onClick => "javascript: document.forms.synchronize.revision.value = ''"), $q->br, "\n";
    print "Specific revision: ",
    $q->textfield(-name => 'revision', -size => 6, -maxlength => 6,
          -onChange => "javascript: document.forms.synchronize.revision_head.checked = 0"), $q->br, $q->br, "\n";

    if ($group =~ m/^Web/i) {
        print $q->checkbox(-name => 'clean_page_template', -checked => 0,
		-value => 1, -label => 'Clean View Cache'), $q->br, "<br/>\n\n\n";

        if ($group eq 'Web Product') {
            print $q->checkbox(-name => 'build_static_page', -checked => 0,
                -value => 1, -label => 'Build Static Home Page'), $q->br, "<br/>\n\n\n";
        }
    }

    if (defined $conf->{update_version_group}) {
         print $q->checkbox(-name => 'update_static_group', -checked => 1,
            -value => $conf->{update_version_group}, -label => 'Update Static Version for CDN'), $q->br, "<br/>\n\n\n";
    }

    if ($conf->{restartable_services}) {
        # restart services
        my $default = 0;
        if ($conf->{conf_type} eq 'maintance' ) {
            print $q->b('As the last step, should Static service be restarted?'), $q->br, "\n";
            print "(The updated code is only effective after restart the Static server)", $q->br, "\n";
        } elsif ($conf->{conf_type} =~ m/^vendor_/i ) {
            print $q->b('As the last step, should vender service be restarted?'), $q->br, "\n";
            print "(Vendor service needs to be restarted if configurations have been changed.)", $q->br, "\n";
        } else {
            print $q->b('As the last step, should Apache be restarted?'), $q->br, "\n";
            print "(Apache needs to be restarted only if Perl modules have been changed in the \"lib\" directory.)", $q->br, "\n";
        }
        print $q->radio_group(-name => 'restart_services', -values => [1, 0],
                              -default => $default, -labels => { 1 => 'Yes', 0 => 'No' }),
        $q->br, $q->br, "\n";
    }

    if ($conf->{reloadable_services}) {
        # reload services
        if ($conf->{conf_type} eq 'named_servers' ) {
            print $q->b('As the last step, should Named be reload?'), $q->br, "\n";
            print "(Named needs to be reloaded if configurations have been changed.)", $q->br, "\n";
        } else {
            print $q->b('As the last step, should unknown service be reload?'), $q->br, "\n";
        }
        print $q->radio_group(-name => 'reload_services', -values => [1, 0], -default => 1,
            -labels => { 1 => 'Yes', 0 => 'No' }),
        $q->br, $q->br, "\n";
    }

    if ($conf->{fast_sync_available}) {
        print $q->checkbox(-name => 'fast_sync', -checked => 1,
                           -value => 1, -label => 'Use fast synchronization'),
            $q->br, "\n";
    }

    print $q->submit(Submit => 'Synchronize servers'), "&nbsp;&nbsp;&nbsp;", $q->reset, "\n";
    print $q->end_form, "\n";

    print $q->hr . $q->a({-href => $q->url }, 'Back to the server type selection page');
    print $q->end_html, "\n";
}

sub sort_servers {
    my ($location_a, $number_a) = ($a =~ /^(.+?)(\d+)$/i);
    my ($location_b, $number_b) = ($b =~ /^(.+?)(\d+)$/i);
    my $location_cmp = $location_a cmp $location_b;
    return $location_cmp ? $location_cmp : $number_a <=> $number_b;
}

sub start_our_job {
    my ($parent_code, $child_code) = @_;

    my $pid = start_job({ stdout_capture => 1,
                          stderr_capture => 1 },
                        '-');
    if ($pid) {
        $parent_code->($pid);
    } else {
        $child_code->();
    }
}

sub watch_our_jobs {
    my ($servers_order, $fields, $code) = @_;

    foreach (keys %status) {
        my $i = position_in_list($servers_order, $fields, $status{$_});
        $servers_order->[$i]->{pid} = $_ if $i >= 0;
    }

    my $next_expected_position = 0;

    while (my ($pid, $event, $data) = watch_jobs()) {
        if ($event eq 'STDOUT') {
            if ($data ne '') {
                $status{$pid}->{STDOUT} .= $data;
            } else {
                $status{$pid}->{STDOUT_done} = 1;
            }
        } elsif ($event eq 'STDERR') {
            if ($data ne '') {
                $status{$pid}->{STDERR} .= $data;
            } else {
                $status{$pid}->{STDERR_done} = 1;
            }
        } elsif ($event eq 'EXIT') {
            $status{$pid}->{EXIT} = $data;
            $status{$pid}->{EXIT_done} = 1;
        }
        if ($status{$pid}->{EXIT_done} and $status{$pid}->{STDOUT_done} and $status{$pid}->{STDERR_done}) {
            my $position = position_in_list($servers_order, $fields, $status{$pid});
            if ($position >= 0) {
                foreach my $i ($next_expected_position .. $#$servers_order) {
                    my $i_pid = $servers_order->[$i]->{pid};
                    if ($status{$i_pid}->{EXIT_done} and $status{$i_pid}->{STDOUT_done} and $status{$i_pid}->{STDERR_done}) {
                        $next_expected_position = $i + 1;
                        $code->($status{$i_pid});
                    } else {
                        last;
                    }
                }
            }
        }
    }
}

sub position_in_list {
    my ($list, $fields, $item) = @_;

    my $i = 0;
    foreach my $list_item (@$list) {
        my $all_fields_match = 1;
        foreach my $field_name (@$fields) {
            if ($list_item->{$field_name} ne $item->{$field_name}) {
                $all_fields_match = 0;
                last;
            }
        }
        if ($all_fields_match) {
            return $i;
        }
        $i++;
    }
    return -1;
}

sub is_ipaddr {
    my ($ip) = @_;

    return $ip =~ m/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/;
}

sub get_ip {
    my ($hostname) = @_;
    if (my $ip = gethostbyname $hostname) {
        return inet_ntoa($ip);
    }
}

