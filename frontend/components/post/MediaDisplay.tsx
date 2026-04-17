"use client";

import Image from "next/image";
import { useState } from "react";
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
    const [isVideoLoading, setIsVideoLoading] = useState(true);

    return (
       <div className="w-full block">
            {(media.media_type == "mp4" || media.media_type == "webm" ||  media.media_type == "mov" ) ?
                <div className="w-full">
                    {isVideoLoading && (
                        <div className="w-full aspect-video rounded-lg bg-zinc-200 animate-pulse" />
                    )}
                    <VideoPlayer
                        className={`w-full overflow-hidden rounded-lg border${isVideoLoading ? " hidden" : ""}`}
                    >
                        <VideoPlayerContent
                        crossOrigin=""
                        muted
                        preload="auto"
                        slot="media"
                        src={media.s3key}
                        onLoadedData={() => setIsVideoLoading(false)}
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
                </div>
            : (media.media_type == "image/jpeg" || media.media_type == "image/png" ) ?
                <Image src={media.s3key} width={200} height={300} alt={"Premium Content Image"} />
            :
                <PDFViewer src={media.s3key} />
            }

       </div>
    );
}