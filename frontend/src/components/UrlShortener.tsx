'use client';

import { useState } from 'react';
import { config } from '@/lib/config';

interface ShortenResponse {
  code: string;
}

interface ErrorResponse {
  error: string;
}

const urlRegex = /^(?:https?:\/\/)?[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*\.[A-Za-z]{2,6}(?::\d{1,5})?(?:[/?#][^\s]*)?$/;

export function UrlShortener() {
  const [url, setUrl] = useState('');
  const [shortenedUrl, setShortenedUrl] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showResult, setShowResult] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setShowResult(false);
    setError('');

    if (!urlRegex.test(url)) {
      setError('Please enter a valid URL (e.g., example.com or https://example.com)');
      return;
    }

    setIsLoading(true);

    try {
      const response = await fetch(`${config.apiBaseUrl}/api/shorten`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      });

      const data = await response.json();

      if (response.ok) {
        const result = data as ShortenResponse;
        setShortenedUrl(`${config.apiBaseUrl}/${result.code}`);
        setShowResult(true);
      } else {
        const errorData = data as ErrorResponse;
        setError(errorData.error || 'Failed to shorten URL');
      }
    } catch (err) {
      setError(`Error: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(shortenedUrl);
      // Could add a toast notification here
    } catch (err) {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = shortenedUrl;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    }
  };

  return (
    <div className="bg-white/95 backdrop-blur-lg rounded-3xl p-12 shadow-2xl max-w-2xl w-full border border-white/30">
      <h1 className="text-center mb-8 text-gray-800 text-4xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
        URL Shortener
      </h1>
      <p className="text-center mb-8 text-gray-600 text-lg">
        Transform long URLs into short, shareable links
      </p>

      <form onSubmit={handleSubmit}>
        <div className="mb-6">
          <label htmlFor="url" className="block mb-2 text-gray-800 font-semibold">
            Original URL
          </label>
          <input
            type="text"
            id="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="w-full p-4 border-2 border-gray-200 rounded-xl text-base transition-all duration-300 bg-white text-gray-900 focus:outline-none focus:border-blue-500 focus:shadow-lg focus:-translate-y-0.5"
            placeholder="example.com/very-long-url-that-needs-shortening"
            required
            suppressHydrationWarning={true}
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full p-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white border-none rounded-xl text-lg font-semibold cursor-pointer transition-all duration-300 shadow-lg hover:-translate-y-0.5 hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? 'Shortening...' : 'Shorten URL'}
        </button>
      </form>

      {showResult && (
        <div className="mt-8 p-6 bg-gradient-to-r from-green-50 to-green-100 rounded-xl border-l-4 border-green-500 animate-in slide-in-from-bottom-4 duration-500">
          <div className="font-semibold text-green-700 mb-2">Your shortened URL:</div>
          <div className="bg-white p-4 rounded-lg border border-green-200 font-mono break-all flex justify-between items-center gap-4">
            <span className="text-gray-900 font-semibold">{shortenedUrl}</span>
            <button
              onClick={copyToClipboard}
              className="bg-green-500 text-white border-none px-4 py-2 rounded-md cursor-pointer text-sm transition-all duration-300 hover:bg-green-600 min-w-[60px]"
            >
              Copy
            </button>
          </div>
        </div>
      )}

      {error && (
        <div className="mt-4 p-4 bg-red-50 rounded-lg border-l-4 border-red-500 text-red-700 animate-in slide-in-from-bottom-4 duration-500">
          {error}
        </div>
      )}
    </div>
  );
}