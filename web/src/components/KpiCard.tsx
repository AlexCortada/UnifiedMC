import type { ReactNode } from 'react';
import { TrendingUp, TrendingDown } from 'lucide-react';

interface KpiCardProps {
  label: string;
  value: string | number;
  change?: number;
  previousValue?: number;
  icon: ReactNode;
  color?: 'blue' | 'green' | 'red' | 'amber' | 'purple';
}

const colorClasses = {
  blue: 'bg-blue-50 text-blue-600 border-blue-200',
  green: 'bg-emerald-50 text-emerald-600 border-emerald-200',
  red: 'bg-red-50 text-red-600 border-red-200',
  amber: 'bg-amber-50 text-amber-600 border-amber-200',
  purple: 'bg-purple-50 text-purple-600 border-purple-200',
};

export function KpiCard({ label, value, change, icon, color = 'blue' }: KpiCardProps) {
  const isPositive = change !== undefined && change >= 0;

  return (
    <div className="card flex items-center gap-4">
      <div className={`flex items-center justify-center w-12 h-12 rounded-lg border ${colorClasses[color]}`}>
        {icon}
      </div>
      <div className="flex-1 min-w-0">
        <p className="kpi-label">{label}</p>
        <p className="kpi-value">{value}</p>
        {change !== undefined && (
          <div className={`flex items-center gap-1 mt-1 ${isPositive ? 'trend-up' : 'trend-down'}`}>
            {isPositive ? <TrendingUp size={14} /> : <TrendingDown size={14} />}
            <span>{isPositive ? '+' : ''}{change.toFixed(1)}% vs last period</span>
          </div>
        )}
      </div>
    </div>
  );
}
