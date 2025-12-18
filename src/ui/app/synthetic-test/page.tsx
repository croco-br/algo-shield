'use client';

import { useState } from 'react';

interface Field {
	name: string;
	type: 'string' | 'number' | 'boolean' | 'date' | 'uuid' | 'email' | 'ip';
	required: boolean;
	format?: string;
}

interface TestResult {
	timestamp: string;
	data: any;
	score: number;
	action: string;
	status: 'success' | 'error';
}

const fieldTypes = [
	{ value: 'string', label: 'String' },
	{ value: 'number', label: 'Number' },
	{ value: 'boolean', label: 'Boolean' },
	{ value: 'date', label: 'Date/Time' },
	{ value: 'uuid', label: 'UUID' },
	{ value: 'email', label: 'Email' },
	{ value: 'ip', label: 'IP Address' },
];

export default function SyntheticTestPage() {
	const [fields, setFields] = useState<Field[]>([
		{ name: 'transaction_id', type: 'uuid', required: true },
		{ name: 'amount', type: 'number', required: true },
		{ name: 'user_email', type: 'email', required: true },
	]);
	const [testResults, setTestResults] = useState<TestResult[]>([]);
	const [loading, setLoading] = useState(false);
	const [shouldStop, setShouldStop] = useState(false);
	const [currentProgress, setCurrentProgress] = useState(0);
	const [numberOfEvents, setNumberOfEvents] = useState(1);
	const [showAddField, setShowAddField] = useState(false);
	const [newField, setNewField] = useState<Partial<Field>>({ name: '', type: 'string', required: true });

	function addField() {
		if (newField.name && newField.type) {
			setFields([...fields, newField as Field]);
			setNewField({ name: '', type: 'string', required: true });
			setShowAddField(false);
		}
	}

	function removeField(index: number) {
		setFields(fields.filter((_, i) => i !== index));
	}

	function generateValue(field: Field): any {
		switch (field.type) {
			case 'uuid':
				return crypto.randomUUID();
			case 'email':
				return `user${Math.floor(Math.random() * 10000)}@example.com`;
			case 'number':
				return Math.floor(Math.random() * 10000) + 100;
			case 'boolean':
				return Math.random() > 0.5;
			case 'date':
				return new Date().toISOString();
			case 'ip':
				return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`;
			case 'string':
			default:
				return `value_${Math.random().toString(36).substring(7)}`;
		}
	}

	function generateSyntheticData(): any {
		const data: any = {};
		fields.forEach((field) => {
			data[field.name] = generateValue(field);
		});
		return data;
	}

	async function runTest() {
		setLoading(true);
		setShouldStop(false);
		setCurrentProgress(0);
		setTestResults([]);

		try {
			for (let i = 0; i < numberOfEvents; i++) {
				if (shouldStop) {
					console.log('Test stopped by user');
					break;
				}

				const syntheticData = generateSyntheticData();

				// Simulate API call
				await new Promise((resolve) => setTimeout(resolve, 300));

				const mockScore = Math.floor(Math.random() * 100);
				const mockAction = mockScore > 80 ? 'block' : mockScore > 50 ? 'review' : 'allow';

				setTestResults((prev) => [
					...prev,
					{
						timestamp: new Date().toISOString(),
						data: syntheticData,
						score: mockScore,
						action: mockAction,
						status: 'success',
					},
				]);

				setCurrentProgress(i + 1);
			}
		} catch (error) {
			console.error('Test failed:', error);
		}

		setLoading(false);
		setCurrentProgress(0);
	}

	function stopTest() {
		setShouldStop(true);
	}

	function exportResults() {
		const dataStr = JSON.stringify(testResults, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `synthetic-test-${Date.now()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}

	function clearResults() {
		setTestResults([]);
	}

	return (
		<div className="max-w-7xl mx-auto px-8">
			<div className="mb-8">
				<h2 className="text-3xl font-semibold mb-2">Synthetic Data Testing</h2>
				<p className="text-gray-500">Generate and test synthetic events with custom field schemas</p>
			</div>

			<div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
				{/* Field Configuration Panel */}
				<div className="bg-white rounded-lg border border-gray-200 p-6">
					<h3 className="text-xl font-semibold mb-2">Event Schema</h3>
					<p className="text-sm text-gray-500 mb-6">Define the fields for your synthetic events</p>

					<div className="space-y-3">
						{fields.map((field, index) => (
							<div key={index} className="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-md">
								<div className="flex-1">
									<div className="font-medium text-gray-900 flex items-center gap-2">
										{field.name}
										{field.required && (
											<span className="px-2 py-0.5 bg-yellow-500 text-white rounded text-xs">Required</span>
										)}
									</div>
									<div className="text-sm text-gray-500">{field.type}</div>
								</div>
								<button
									onClick={() => removeField(index)}
									className="text-red-600 text-xl leading-none hover:opacity-80"
									title="Remove field"
								>
									✕
								</button>
							</div>
						))}

						{showAddField ? (
							<div className="p-4 bg-gray-50 border-2 border-dashed border-gray-300 rounded-md space-y-3">
								<input
									type="text"
									value={newField.name || ''}
									onChange={(e) => setNewField({ ...newField, name: e.target.value })}
									placeholder="Field name"
									className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								/>
								<select
									value={newField.type || 'string'}
									onChange={(e) => setNewField({ ...newField, type: e.target.value as Field['type'] })}
									className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								>
									{fieldTypes.map((type) => (
										<option key={type.value} value={type.value}>
											{type.label}
										</option>
									))}
								</select>
								<label className="flex items-center gap-2 text-sm">
									<input
										type="checkbox"
										checked={newField.required || false}
										onChange={(e) => setNewField({ ...newField, required: e.target.checked })}
									/>
									Required
								</label>
								<div className="flex gap-2">
									<button
										onClick={addField}
										className="px-4 py-2 bg-indigo-600 text-white rounded text-sm hover:bg-indigo-700 transition-colors"
									>
										Add
									</button>
									<button
										onClick={() => setShowAddField(false)}
										className="px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
									>
										Cancel
									</button>
								</div>
							</div>
						) : (
							<button
								onClick={() => setShowAddField(true)}
								className="w-full px-4 py-2 border border-dashed border-gray-300 rounded text-sm text-gray-500 hover:border-indigo-500 hover:text-indigo-600 transition-colors"
							>
								+ Add Field
							</button>
						)}
					</div>
				</div>

				{/* Test Configuration Panel */}
				<div className="bg-white rounded-lg border border-gray-200 p-6">
					<h3 className="text-xl font-semibold mb-2">Test Configuration</h3>
					<p className="text-sm text-gray-500 mb-6">Configure and run your tests</p>

					<div className="mb-6">
						<label htmlFor="numEvents" className="block text-sm font-medium text-gray-700 mb-2">
							Number of Events
						</label>
						<input
							id="numEvents"
							type="number"
							min="1"
							max="1000"
							value={numberOfEvents}
							onChange={(e) => setNumberOfEvents(parseInt(e.target.value) || 1)}
							disabled={loading}
							className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-60"
						/>
					</div>

					{loading ? (
						<>
							<div className="mb-4 p-4 bg-gray-50 border border-gray-200 rounded-md">
								<div className="flex justify-between items-center mb-3 text-sm">
									<span>Processing events...</span>
									<span className="font-semibold text-indigo-600">
										{currentProgress} / {numberOfEvents}
									</span>
								</div>
								<div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
									<div
										className="h-full bg-gradient-to-r from-indigo-600 to-purple-600 transition-all duration-300"
										style={{ width: `${(currentProgress / numberOfEvents) * 100}%` }}
									></div>
								</div>
							</div>

							<button
								onClick={stopTest}
								className="w-full py-3 px-4 bg-red-600 text-white rounded-md font-medium hover:bg-red-700 transition-colors"
							>
								⏹ Stop Test
							</button>
						</>
					) : (
						<button
							onClick={runTest}
							disabled={fields.length === 0}
							className="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
						>
							▶ Run Test
						</button>
					)}

					{testResults.length > 0 && (
						<>
							<div className="grid grid-cols-2 gap-4 mt-6 p-6 bg-gray-50 rounded-md">
								<div className="text-center">
									<div className="text-3xl font-bold text-indigo-600">{testResults.length}</div>
									<div className="text-sm text-gray-500 mt-1">Events Tested</div>
								</div>
								<div className="text-center">
									<div className="text-3xl font-bold text-indigo-600">
										{testResults.filter((r) => r.action === 'block').length}
									</div>
									<div className="text-sm text-gray-500 mt-1">Blocked</div>
								</div>
								<div className="text-center">
									<div className="text-3xl font-bold text-indigo-600">
										{testResults.filter((r) => r.action === 'review').length}
									</div>
									<div className="text-sm text-gray-500 mt-1">Review</div>
								</div>
								<div className="text-center">
									<div className="text-3xl font-bold text-indigo-600">
										{testResults.filter((r) => r.action === 'allow').length}
									</div>
									<div className="text-sm text-gray-500 mt-1">Allowed</div>
								</div>
							</div>

							<div className="flex gap-2 mt-4">
								<button
									onClick={exportResults}
									className="flex-1 px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
								>
									Export Results
								</button>
								<button
									onClick={clearResults}
									className="flex-1 px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
								>
									Clear
								</button>
							</div>
						</>
					)}
				</div>
			</div>

			{/* Results Panel */}
			{testResults.length > 0 && (
				<div className="bg-white rounded-lg border border-gray-200 p-6">
					<h3 className="text-xl font-semibold mb-6">Test Results</h3>

					<div className="space-y-4">
						{testResults.map((result, index) => (
							<div key={index} className="border border-gray-200 rounded-md p-4">
								<div className="flex items-center gap-4 mb-2">
									<span className="font-bold text-gray-500">#{index + 1}</span>
									<span className="text-sm text-gray-500">
										{new Date(result.timestamp).toLocaleTimeString()}
									</span>
									<span
										className={`px-3 py-1 rounded-full text-xs font-medium ${
											result.action === 'block'
												? 'bg-red-500 text-white'
												: result.action === 'review'
													? 'bg-yellow-500 text-white'
													: 'bg-green-500 text-white'
										}`}
									>
										{result.action}
									</span>
									<span className="ml-auto font-medium">Score: {result.score}</span>
								</div>
								<details className="mt-3">
									<summary className="cursor-pointer text-sm text-indigo-600 hover:text-indigo-700 select-none">
										View Event Data
									</summary>
									<pre className="mt-3 p-4 bg-gray-50 rounded border border-gray-200 overflow-x-auto text-sm">
										{JSON.stringify(result.data, null, 2)}
									</pre>
								</details>
							</div>
						))}
					</div>
				</div>
			)}
		</div>
	);
}
