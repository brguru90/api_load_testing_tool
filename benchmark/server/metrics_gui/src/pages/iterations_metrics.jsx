import React from 'react'
import { useLocation } from 'react-router-dom';
import APIPerSecondMetrics from '../components/each_api_metrics/each_iteration_in_detail';

export default function IterationsMetrics(props) {
  const { search } = useLocation();
  const APIindex = new URLSearchParams(search).get("api_index")
  return <APIPerSecondMetrics APIindex={APIindex} />
}
