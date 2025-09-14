'use client';

export function MetricsDashboard() {
  return (
    <div className="bg-white/95 backdrop-blur-lg rounded-3xl p-8 shadow-2xl w-full max-w-7xl border border-white/30">
      <h2 className="text-center mb-8 text-gray-800 text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
        Live Metrics Dashboard
      </h2>

      {/* Time Series Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
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
        <div className="bg-white rounded-xl p-4 shadow-lg border border-gray-200">
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