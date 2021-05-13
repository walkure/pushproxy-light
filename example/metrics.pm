use 5.000;
use strict;
use warnings;

package _metric_value;

sub new {
	my $class = shift;
	my $value = shift;
	my $self = {'value' => ''.$value };
	bless $self, $class;
	return $self;
}

sub add_label {
	my ($self,$key,$value) = @_;
	unless(defined $self->{'labels'}){
		$self->{'labels'} = {};
	}
	
	$self->{'labels'}->{''.$key} = ''.$value;
	
	$self;
}

sub damn_object {
	my $self = shift;
	
	if(defined $self->{'labels'}){
		{'value' => $self->{'value'},'labels' => $self->{'labels'}};
	}else{
		{'value' => $self->{'value'}};
	}
}

package _metric_item;

sub new {
	my ($class,%args) = @_;
	
	my $self = {
		'type' => $args{'type'},
		'help' => $args{'help'},
		'metrics' => []
	};
	
	bless $self, $class;
	return $self;
}


sub add_value {
	my($self,$value) = @_;
	
	my $metric = _metric_value->new($value);
	push(@{$self->{'metrics'}},$metric);
	$metric;
}

sub damn_object {
	my $self = shift;
	
	my @damned_metrics;
	foreach my $metric(@{$self->{'metrics'}}){
		push(@damned_metrics,$metric->damn_object());
	}
	
	{
		'help'    => $self->{'help'},
		'type'    => $self->{'type'},
		'metrics' => \@damned_metrics,
	};
}

package metrics;
use JSON 'encode_json';

sub new {
	my $class = shift;

	my $self = {};
	$self->{'metrics'} = {};
	bless $self, $class;
	return $self;
}

sub add_metric {
	my($self,%args) = @_;
	
	$self->{'metrics'}->{$args{'name'}} = _metric_item->new(%args)
}

sub to_json {
	my $self = shift;
	my $damned = {};
	
	foreach my $key (keys %{$self->{'metrics'}}) {
		$damned->{$key} = $self->{'metrics'}->{$key}->damn_object();
	}
	
	encode_json($damned);
}

1;

