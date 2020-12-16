export default({}, inject) => {
	const server = location.origin !== "http://localhost:3000" ? location.origin : "http://localhost:7473"
	inject('server', server)
}