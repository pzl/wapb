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
		</v-col>
		<v-col cols="12" md="">
			<file-row v-for="f in files" :key="f.id" v-bind="f" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import FileRow from '~/components/fileRow'

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
			files: []
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
				let form = new FormData();
				for (const f of this.toCreate.files) {
					form.append(f.name, f)
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
	components: { CreateMeta, FileRow }

}
</script>