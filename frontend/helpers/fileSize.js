/**
 * Bytes as human readable text
 */
export default (bytes, si=false, dp=1) => {
	const thresh = si ? 1000 : 1024;

	if (bytes < thresh) {
		return `${bytes}B`
	}

	const units = si
		? ['kB', 'MB', 'GB', 'TB', 'PB']
		: ['KiB', 'MiB', 'GiB', 'TiB', 'PiB'];

	let i = -1;
	const r = 10 ** dp;
	do {
		bytes /= thresh;
		i++;
	} while (Math.round(bytes * r) / r >= thresh && i < units.length - 1);

	return `${bytes.toFixed(dp)} ${units[i]}`;
}