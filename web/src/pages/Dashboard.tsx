import { useEffect, useState } from 'react';
import {
  Monitor,
  Wifi,
  WifiOff,
  ShieldCheck,
  ShieldAlert,
  AlertTriangle,
  FileText,
  Clock,
  UserPlus,
  UserMinus,
  RefreshCw,
} from 'lucide-react';
import { KpiCard } from './KpiCard';
import { fetchDashboardSummary, fetchDevices, type DashboardSummary, type Device } from '../services/api';

export function Dashboard() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastRefresh, setLastRefresh] = useState<Date>(new Date());

  const loadData = async () => {
    try {
      const [summaryData, devicesData] = await Promise.all([
        fetchDashboardSummary(),
        fetchDevices(),
      ]);
      setSummary(summaryData);
      setDevices(devicesData.items);
      setLastRefresh(new Date());
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 60000); // Auto-refresh every 60s
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="flex items-center gap-3 text-gray-500">
          <RefreshCw className="animate-spin" size={24} />
          <span className="text-lg">Loading dashboard...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <AlertTriangle className="mx-auto text-red-500 mb-4" size={48} />
          <h2 className="text-xl font-semibold text-gray-900 mb-2">Connection Error</h2>
          <p className="text-gray-500 mb-4">{error}</p>
          <button
            onClick={loadData}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">IT Operations Portal</h1>
            <p className="text-sm text-gray-500">Unified device management & operations</p>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-xs text-gray-400">
              Last updated: {lastRefresh.toLocaleTimeString()}
            </span>
            <button
              onClick={loadData}
              className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition"
              title="Refresh"
            >
              <RefreshCw size={18} />
            </button>
          </div>
        </div>
      </header>

      <main className="p-6 max-w-7xl mx-auto">
        {/* KPI Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
          <KpiCard
            label="Total Devices"
            value={summary?.total_devices.value ?? 0}
            change={summary?.total_devices.change_percent}
            icon={<Monitor size={22} />}
            color="blue"
          />
          <KpiCard
            label="Online"
            value={summary?.online_devices.value ?? 0}
            icon={<Wifi size={22} />}
            color="green"
          />
          <KpiCard
            label="Offline"
            value={summary?.offline_devices.value ?? 0}
            icon={<WifiOff size={22} />}
            color="red"
          />
          <KpiCard
            label="Compliance"
            value={`${summary?.compliance_rate.rate ?? 0}%`}
            icon={<ShieldCheck size={22} />}
            color="purple"
          />
          <KpiCard
            label="Critical Vulns"
            value={summary?.critical_vulnerabilities.total ?? 0}
            icon={<ShieldAlert size={22} />}
            color="amber"
          />
        </div>

        {/* Secondary KPI Row */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
          <KpiCard
            label="Open Incidents"
            value={summary?.open_incidents.total ?? 0}
            icon={<FileText size={22} />}
            color="blue"
          />
          <KpiCard
            label="Patch Rate"
            value={`${summary?.patch_compliance.overall_rate ?? 0}%`}
            icon={<ShieldCheck size={22} />}
            color="green"
          />
          <KpiCard
            label="SLA Breaches"
            value={summary?.sla_breaches.active ?? 0}
            icon={<Clock size={22} />}
            color="red"
          />
          <KpiCard
            label="New Hires"
            value={summary?.new_hires.total ?? 0}
            icon={<UserPlus size={22} />}
            color="blue"
          />
          <KpiCard
            label="Terminations"
            value={summary?.terminations.total ?? 0}
            icon={<UserMinus size={22} />}
            color="amber"
          />
        </div>

        {/* Device List */}
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">Devices</h2>
            <span className="text-sm text-gray-500">{devices.length} of {summary?.total_devices.value ?? 0} shown</span>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Device</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">OS</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Type</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Status</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Compliance</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Source</th>
                  <th className="text-left py-3 px-2 font-medium text-gray-500">Last Seen</th>
                </tr>
              </thead>
              <tbody>
                {devices.map((device) => (
                  <tr key={device.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-3 px-2">
                      <div>
                        <p className="font-medium text-gray-900">{device.canonical_name}</p>
                        <p className="text-xs text-gray-400">{device.serial_number || device.external_id}</p>
                      </div>
                    </td>
                    <td className="py-3 px-2 capitalize">{device.os_type} {device.os_version}</td>
                    <td className="py-3 px-2 capitalize">{device.asset_type}</td>
                    <td className="py-3 px-2">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                        device.status === 'active'
                          ? 'bg-emerald-100 text-emerald-700'
                          : 'bg-gray-100 text-gray-600'
                      }`}>
                        {device.status}
                      </span>
                    </td>
                    <td className="py-3 px-2">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                        device.compliance_status === 'compliant'
                          ? 'bg-emerald-100 text-emerald-700'
                          : device.compliance_status === 'non_compliant'
                          ? 'bg-red-100 text-red-700'
                          : 'bg-gray-100 text-gray-600'
                      }`}>
                        {device.compliance_status || 'unknown'}
                      </span>
                    </td>
                    <td className="py-3 px-2">
                      <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700 capitalize">
                        {device.connector_type.replace('_', ' ')}
                      </span>
                    </td>
                    <td className="py-3 px-2 text-gray-500">
                      {new Date(device.last_seen).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}
