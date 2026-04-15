"use client";

import Image from "next/image";
import {
  VideoPlayer,
  VideoPlayerContent,
  VideoPlayerControlBar,
  VideoPlayerMuteButton,
  VideoPlayerPlayButton,
  VideoPlayerSeekBackwardButton,
  VideoPlayerSeekForwardButton,
  VideoPlayerTimeDisplay,
  VideoPlayerTimeRange,
  VideoPlayerVolumeRange,
} from "@/components/kibo-ui/video-player";

import dynamic from "next/dynamic";
import { Media } from "@/api";

const PDFViewer = dynamic(() => import("./PDF"), { ssr: false });

type MediaDisplayProps = React.ComponentProps<"div"> & {
    media: Media;
};

export default function MediaDisplay({ media, className, ...props }: MediaDisplayProps) {

    return (
       <div className="w-full block">
            {(media.media_type == "mp4" || media.media_type == "webm" ||  media.media_type == "mov" ) ?
                <VideoPlayer className="overflow-hidden rounded-lg border">
                    <VideoPlayerContent
                    crossOrigin=""
                    muted
                    preload="auto"
                    slot="media"
                    src={media.s3key}
                    />
                    <VideoPlayerControlBar>
                    <VideoPlayerPlayButton />
                    <VideoPlayerSeekBackwardButton />
                    <VideoPlayerSeekForwardButton />
                    <VideoPlayerTimeRange />
                    <VideoPlayerTimeDisplay showDuration />
                    <VideoPlayerMuteButton />
                    <VideoPlayerVolumeRange />
                    </VideoPlayerControlBar>
                </VideoPlayer>
            : (media.media_type == "jpeg" || media.media_type == "png" ) ? 
                <Image src={media.s3key} width={200} height={300} alt={"Premium Content Image"} />
            :
                <PDFViewer src={media.s3key} />
            }

       </div>
    );
}