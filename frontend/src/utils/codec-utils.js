export default class CodecUtil {
	static isVideoSupported = (videoCodec) => {
		return videoCodec == 'h264';
	};

	static isAudioSupported = (audioCodec) => {
		return audioCodec == 'aac';
	};
}
