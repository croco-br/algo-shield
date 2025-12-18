'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/stores/auth';
import { api } from '@/lib/api';

export default function LoginPage() {
	const { setToken } = useAuth();
	const router = useRouter();
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [name, setName] = useState('');
	const [activeTab, setActiveTab] = useState<'login' | 'register'>('login');
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState('');

	async function handleLogin() {
		if (!email || !password) {
			setError('Please enter email and password');
			return;
		}

		setLoading(true);
		setError('');

		try {
			const response = await api.post<{ token: string; user: any }>('/api/v1/auth/login', {
				email,
				password,
			});

			await setToken(response.token);
			router.push('/');
		} catch (e: any) {
			setError(e.message || 'Login failed. Please try again.');
		} finally {
			setLoading(false);
		}
	}

	async function handleRegister() {
		if (!email || !password || !name) {
			setError('Please fill in all fields');
			return;
		}

		if (password.length < 6) {
			setError('Password must be at least 6 characters');
			return;
		}

		setLoading(true);
		setError('');

		try {
			const response = await api.post<{ token: string; user: any }>('/api/v1/auth/register', {
				email,
				password,
				name,
			});

			await setToken(response.token);
			router.push('/');
		} catch (e: any) {
			setError(e.message || 'Registration failed. Please try again.');
		} finally {
			setLoading(false);
		}
	}

	return (
		<div className="min-h-screen flex items-center justify-center bg-purple-600 p-8">
			<div className="bg-white rounded-xl shadow-xl w-full max-w-md p-10">
				<div className="mb-8">
					<div className="flex items-center gap-3 mb-2">
						<img src="/gopher.png" alt="AlgoShield" className="w-12 h-12 object-contain" />
						<h1 className="text-3xl font-semibold text-gray-900">AlgoShield</h1>
					</div>
					<p className="text-sm text-gray-500">Fraud Detection & Anti-Money Laundering</p>
				</div>

				<div className="flex gap-2 mb-6 border-b-2 border-gray-200">
					<button
						onClick={() => {
							setActiveTab('login');
							setError('');
						}}
						className={`flex-1 py-3 text-sm font-medium transition-all border-b-2 -mb-[2px] ${
							activeTab === 'login'
								? 'text-indigo-600 border-indigo-600 font-semibold'
								: 'text-gray-500 border-transparent'
						}`}
					>
						Login
					</button>
					<button
						onClick={() => {
							setActiveTab('register');
							setError('');
						}}
						className={`flex-1 py-3 text-sm font-medium transition-all border-b-2 -mb-[2px] ${
							activeTab === 'register'
								? 'text-indigo-600 border-indigo-600 font-semibold'
								: 'text-gray-500 border-transparent'
						}`}
					>
						Register
					</button>
				</div>

				{error && (
					<div className="bg-red-50 border border-red-200 rounded-md p-3 mb-6 text-sm text-red-600">
						{error}
					</div>
				)}

				{activeTab === 'login' ? (
					<form
						onSubmit={(e) => {
							e.preventDefault();
							handleLogin();
						}}
						className="space-y-6"
					>
						<div>
							<label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
								Email
							</label>
							<input
								id="email"
								type="email"
								value={email}
								onChange={(e) => setEmail(e.target.value)}
								placeholder="user@example.com"
								required
								disabled={loading}
								className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
							/>
						</div>

						<div>
							<label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
								Password
							</label>
							<input
								id="password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								placeholder="••••••••"
								required
								disabled={loading}
								className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
							/>
						</div>

						<button
							type="submit"
							disabled={loading}
							className="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
						>
							{loading ? 'Signing in...' : 'Sign In'}
						</button>
					</form>
				) : (
					<form
						onSubmit={(e) => {
							e.preventDefault();
							handleRegister();
						}}
						className="space-y-6"
					>
						<div>
							<label htmlFor="reg-name" className="block text-sm font-medium text-gray-700 mb-2">
								Name
							</label>
							<input
								id="reg-name"
								type="text"
								value={name}
								onChange={(e) => setName(e.target.value)}
								placeholder="Your Name"
								required
								disabled={loading}
								className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
							/>
						</div>

						<div>
							<label htmlFor="reg-email" className="block text-sm font-medium text-gray-700 mb-2">
								Email
							</label>
							<input
								id="reg-email"
								type="email"
								value={email}
								onChange={(e) => setEmail(e.target.value)}
								placeholder="user@example.com"
								required
								disabled={loading}
								className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
							/>
						</div>

						<div>
							<label htmlFor="reg-password" className="block text-sm font-medium text-gray-700 mb-2">
								Password
							</label>
							<input
								id="reg-password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								placeholder="••••••••"
								required
								minLength={6}
								disabled={loading}
								className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
							/>
							<small className="text-xs text-gray-500 mt-1 block">Minimum 6 characters</small>
						</div>

						<button
							type="submit"
							disabled={loading}
							className="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
						>
							{loading ? 'Creating account...' : 'Create Account'}
						</button>
					</form>
				)}
			</div>
		</div>
	);
}
