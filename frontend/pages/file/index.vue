<template>
	<v-row>
		<v-col>
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-file-input counter multiple show-size chips truncate-length="20"
				name="file" label="File" hint="File"
				v-model="toCreate.file"
			/>
			<create-meta
				ref="meta"
				v-bind.sync="toCreate"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="!!toCreate.file">Create</v-btn>
		</v-col>
		<v-col>
			<file-row v-for="f in files" :key="f.id" v-bind="f" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import FileRow from '~/components/fileRow'

const createFactory = () => {
	return {
		file: null,
		burn: false,
		ttl: null,
		hidden: false,
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
	async asyncData({ $http }) {
		const files = await $http.$get('http://localhost:7473/api/v1/file').then(d => d.data)
		return { files }
	},
	methods: {
		async create() {
			if (!this.toCreate.file) {
				this.err = { message: 'Cannot create empty file' };
				return
			}
			this.alert = null;
			this.isLoading = true;
			const data = await this.$http.$post('http://localhost:7473/api/v1/file', this.toCreate).then(d => {
				this.isLoading = false;
				this.toCreate = createFactory()
				this.$refs.meta.hideTTL()
				this.files.push(d)
				this.alert = {
					type: "success",
					message: `Created ${d.id}`,
				}
			}).catch(e => {
				this.isLoading = false;
				this.alert = {
					type: "error",
					message: e.message,
				}
			})
		},
	},
	components: { CreateMeta, FileRow }

}
</script>