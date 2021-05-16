use strict;
use warnings;

use HTTP::Tiny;
use Digest::SHA 'sha256_hex';
use Time::HiRes qw(time);

use lib qw(./);
use metrics;

my $m = new metrics;
my $ts = int(time() * 1000);
my $it = $m->add_metric('name'=>'absolute_humidity','help'=>'Absolute Humidity g/m^3','type' => 'gauge');
$it->add_value(18.2)->add_label('place','inside')->add_timestamp($ts);
$it->add_value(19.3)->add_label('place','outside');

$it = $m->add_metric('name'=>'relative_humidity','help'=>'Relative Humidity percent','type' => 'gauge');
$it->add_value(44.2)->add_label('place','inside');
$it->add_value(76.5)->add_label('place','outside');

$it = $m->add_metric('name'=>'temperature','help'=>'Temperature in Celsius degree','type' => 'gauge');
$it->add_value(24.2)->add_label('place','inside');
$it->add_value(36.56)->add_label('place','outside');

$it = $m->add_metric('name'=>'pressure','help'=>'Pressure hPa','type' => 'gauge');
$it->add_value(24.2,$ts);

my $key = 'foobarfoobarfoobarfoobarfoobarfoobarfoobar';

my $body = $m->to_json();


my $res = HTTP::Tiny->new->post_form('https://po:yo@example.jp/push/hoge:8080',{
	'body' => $body,
	'signature' => '5'.sha256_hex($body.$key),
#my $res = HTTP::Tiny->new->post('https://po:yo@example.jp/push/hoge:8080',{
#	'content' => $body,
#	'headers' => {
#		'x-signature' => '5'.sha256_hex($body.$key),
#		'content-type' => 'application/json',
}) or die "Unable to get document: $!";


print "req: $res->{status}\n";
print $res->{content};
