'use client';

export function MetricsDashboard() {
  return (
    <div className="bg-white/10 backdrop-blur-xl rounded-2xl p-8 shadow-2xl w-full max-w-7xl border border-white/20">
      <div className="text-center mb-8">
        <div className="flex items-center justify-center gap-3 mb-4">
          <div className="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-xl flex items-center justify-center">
            <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
        <h2 className="text-3xl font-bold text-white mb-2">
          Live Metrics Dashboard
        </h2>
        <p className="text-white/70">
          Real-time analytics and performance insights
        </p>
      </div>

      {/* Time Series Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white/10 backdrop-blur-sm rounded-xl p-4 shadow-lg border border-white/20">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=1&refresh=5s"
            className="w-full h-72 border-none rounded-lg"
            title="Metrics Panel 1"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=2&refresh=5s"
            className="w-full h-72 border-none rounded-lg"
            title="Metrics Panel 2"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=3&refresh=5s"
            className="w-full h-72 border-none rounded-lg"
            title="Metrics Panel 3"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=4&refresh=5s"
            className="w-full h-72 border-none rounded-lg"
            title="Metrics Panel 4"
          />
        </div>
      </div>

      {/* Statistics Panels */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white/10 backdrop-blur-sm rounded-xl p-4 shadow-lg border border-white/20">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=5&refresh=5s"
            className="w-full h-48 border-none rounded-lg"
            title="Stat Panel 1"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=6&refresh=5s"
            className="w-full h-48 border-none rounded-lg"
            title="Stat Panel 2"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=7&refresh=5s"
            className="w-full h-48 border-none rounded-lg"
            title="Stat Panel 3"
          />
        </div>

        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
          <iframe
            src="http://localhost:3000/d-solo/ac5564a1-fbfc-4f71-9e4a-ae1cb25025f1/metrics?orgId=1&from=now-1h&to=now&timezone=browser&panelId=8&refresh=5s"
            className="w-full h-48 border-none rounded-lg"
            title="Stat Panel 4"
          />
        </div>
      </div>
    </div>
  );
}