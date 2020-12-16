<template>
	<v-row>
		<v-col cols="12" md="">
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-file-input counter multiple show-size chips truncate-length="20"
				name="file" label="File" hint="File"
				v-model="toCreate.files"
			/>
			<create-meta
				ref="meta"
				v-bind.sync="toCreate.meta"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.files.length == 0">Create</v-btn>
			<instructions v-bind="instructions" />
		</v-col>
		<v-col cols="12" md="">
			<file-row v-for="f in files" :key="f.id" v-bind="f" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import FileRow from '~/components/fileRow'
import Instructions from '~/components/instructions'

const createFactory = () => {
	return {
		meta: {
			burn: false,
			ttl: null,
			hidden: false,			
		},
		files: [],
	}
} 

export default {
	data () {
		return {
			toCreate: createFactory(),
			isLoading: false,
			alert: null,
			files: [],

			instructions: {
				notes: `<p>Uploading files required 2 API requests. First, to create a file ID. This has no required fields, just make a POST request to ${this.$server}/api/v1/file. You specify the optional fields (burn, ttl, hidden) during this request.</p>
						<pre>curl -X POST ${this.$server}/api/v1/file</pre>
						<p>This will return a file id (we will call this GID for group-id). Now you upload files to ${this.$server}/api/v1/file/GID. Instructions for doing this file upload are below</p>`,
				fields: [],
				tools: [
					{
						name: "Curl",
						actions: [
							{
								name: "Create",
								methods: [
									{ name: "Single File", code: `curl -F nickname@YOURFILE ${ this.$server }/api/v1/file/GID` },
									{ name: "Multiple Files", code: `curl -F nickname@path/to/file -F file2@path/to/another ${ this.$server }/api/v1/file/GID` },
								]
							},
							{
								name: "Retrieve",
								methods: [
									{ name: "File Group Info", code: `curl ${ this.$server }/api/v1/file/GID` },
									{ name: "File Contents (to terminal output)", code: `curl ${ this.$server }/api/v1/file/GID/FID` },
									{ name: "File Contents (to file)", code: `curl ${ this.$server }/api/v1/file/GID/FID -o myfile.name` },
									{ name: "(text) File Contents into line counter", code: `curl -s ${ this.$server }/api/v1/file/GID/FID | wc -l` },
									{ name: "(image) File Contents into resizer", code: `curl -s ${ this.$server }/api/v1/file/GID/FID | convert - -resize 800x800 output.png` },
								]
							},
							{ name: "Delete", methods: [{ code: `curl -X DELETE ${ this.$server }/api/v1/file/GID` }] },
						]
					},
					{
						name: "wget",
						actions: [
							{
								name: "Retrieve",
								notes: `You will need an ID for a file upload group (GID), and from the group info, an individual file ID (FID)`,
								methods: [
									{
										name: "File Contents",
										code: `wget ${ this.$server }/api/v1/text/GID/FID`,
									}
								]
							}
						]
					},
					{
						name: "HTTPie",
						actions: [
							{
								name: "Create",
								methods: [
									{ name: "Upload one", code: `http -f POST ${ this.$server }/api/v1/file/GID anyName@FILENAME` },
									{ name: "Multiple Files", code: `http -f POST ${ this.$server }/api/v1/file/GID file1@path/to/file file2@path/to/another` },
								]
							},
							{
								name: "Retrieve",
								methods: [
									{ name: "File Group Info", code: `http ${ this.$server }/api/v1/file/GID` },
									{ name: "File Contents", code: `http ${ this.$server }/api/v1/file/GID/FID` },
								]
							},
							{ name: "Delete", methods: [{ code: `http DELETE ${ this.$server }/api/v1/file/GID` }] },
						]
					}
				]
			}
		}
	},
	async asyncData({ $http, $server }) {
		const files = await $http.$get(`${$server}/api/v1/file`).then(d => d.data)
		return { files }
	},
	methods: {
		async create() {
			if (!this.toCreate.files || this.toCreate.files.length == 0) {
				this.err = { message: 'Cannot create empty file' };
				return
			}
			this.alert = null;
			this.isLoading = true;
			const data = await this.$http.$post(`${this.$server}/api/v1/file`, this.toCreate.meta).catch(e => {
				this.isLoading = false;
				this.alert = {
					type: "error",
					message: `Error creating File Entry: ${e.message}`,
				}
			})
			if (data) {
				data.files = [];
				let form = new FormData();
				for (const f of this.toCreate.files) {
					form.append(f.name, f)
					data.files.push({
						name: f.name,
						filename: f.name,
						size: f.size,
						id: '',
						mime: f.type,
					})
				}
				const uploads = await this.$http.$post(`${this.$server}/api/v1/file/${data.id}`, form)
					.then(d => {
						this.isLoading = false;
						this.toCreate = createFactory()
						this.$refs.meta.hideTTL()
						this.files.push(data)
						this.alert = {
							type: "success",
							message: `Created ${data.id}`,
						}
					}).catch(e => {
						this.isLoading = false;
						this.alert = {
							type: "error",
							message: `Error uploading files: ${e.message}`
						}
					})
			}
		},
	},
	components: { CreateMeta, FileRow, Instructions }

}
</script>