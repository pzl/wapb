<template>
	<v-card>
		<div class="d-flex flex-no-wrap justify-space-between">

			<div>
				<v-card-title class="headline">{{ name }}</v-card-title>
				<v-card-subtitle>{{ filename }}, {{ formattedSize }}</v-card-subtitle>
				<v-card-text ></v-card-text>
				<v-card-actions>
					<v-btn class="ml-2 mt-3" fab icon right :href="link" target="_blank">
						<v-icon :title="mime">{{ icon }}</v-icon>
					</v-btn>
				</v-card-actions>
			</div>
			<v-avatar v-if="hasPreview" class="ma-3" size="125" tile>
				<v-img v-if="!burn" :src="link"></v-img>
				<v-sheet v-else color="grey lighten-2" height="100" width="100" shaped><v-icon color="red" x-large>mdi-fire</v-icon></v-sheet>
			</v-avatar>
		</div>
	</v-card>
</template>


<script>
import fileSize from '~/helpers/fileSize'


export default {
	props: {
		group_id: {},
		id: {},
		name: {},
		filename: {},
		size: {},
		mime: {},
		burn: {},
	},
	data() {
		return {

		}
	},
	computed: {
		formattedSize() {
			return fileSize(this.size || 0)
		},
		icon() {
			const mimetype = this.mime.split(";")[0];
			const parts = mimetype.split("/");

			switch(parts[0]) {
				case "image": return 'mdi-file-image';
				case "audio": return 'mdi-file-music-outline';
				case "video": return 'mdi-file-video-outline';
			}


			switch (mimetype) {
				case 'text/calendar': return 'mdi-'
				case "text/csv":
					return 'mdi-file-table'
				case "text/html":
				case "text/css":
				case "application/x-csh":
				case "application/x-sh":
				case "text/javascript":
				case "application/json":
				case "application/x-httpd-php":
				case "application/xml":
				case "text/xml":
					return 'mdi-file-code-outline';
				case "application/pdf":
					return 'mdi-file-pdf-outline';
				case "text/plain":
				case "application/rtf":
					return 'mdi-file-document-outline';
				case "application/msword":
				case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
					return 'mdi-file-word-outline';
				case "application/vnd.ms-powerpoint":
				case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
					return 'mdi-file-chart-outline';
				case "application/vnd.ms-excel":
				case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
					return 'mdi-file-excel-outline';
				case "application/x-bzip":
				case "application/x-bzip2":
				case "application/gzip":
				case "application/vnd.rar":
				case "application/x-tar":
				case "application/zip":
				case "application/x-7z-compressed":
					return 'mdi-zip-box'
			}

			return 'mdi-file'
		},
		link() {
			return `/api/v1/file/${this.group_id}/${this.id}`
		},
		hasPreview() {
			return this.burn || this.mime.split(";")[0].split("/")[0] == "image" && this.size < 4*2**20;
		}
	}
}
</script>

<style>

</style>