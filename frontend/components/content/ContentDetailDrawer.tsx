
import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getContentDetail } from '../../services/apiClient';
import { ContentDetail, ContentMetrics, ContentScores, ContentType } from '../../types';
import { Drawer } from '../ui/Drawer';
import { Badge } from '../ui/Badge';
import { formatDateTime } from '../../utils/dateFormatter';

interface ContentDetailDrawerProps {
  contentId: number | null;
  onClose: () => void;
}

const DetailItem: React.FC<{ label: string; value: React.ReactNode }> = ({ label, value }) => (
    <div className="flex flex-col">
        <span className="text-sm text-gray-500">{label}</span>
        <span className="text-lg font-semibold text-gray-800">{value}</span>
    </div>
);

const MetricsDisplay: React.FC<{ metrics: ContentMetrics, type: ContentType }> = ({ metrics, type }) => (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {type === 'video' && <>
            <DetailItem label="Views" value={metrics.views?.toLocaleString() ?? 'N/A'} />
            <DetailItem label="Likes" value={metrics.likes?.toLocaleString() ?? 'N/A'} />
            <DetailItem label="Duration (s)" value={metrics.duration_seconds ?? 'N/A'} />
        </>}
        {type === 'article' && <>
            <DetailItem label="Reading Time (min)" value={metrics.reading_time_minutes ?? 'N/A'} />
            <DetailItem label="Reactions" value={metrics.reactions?.toLocaleString() ?? 'N/A'} />
            <DetailItem label="Comments" value={metrics.comments?.toLocaleString() ?? 'N/A'} />
        </>}
    </div>
);

const ScoresDisplay: React.FC<{ scores: ContentScores }> = ({ scores }) => (
    <div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <DetailItem label="Base" value={scores.base} />
            <DetailItem label="Type Coeff." value={scores.type_coefficient} />
            <DetailItem label="Freshness" value={scores.freshness} />
            <DetailItem label="Engagement" value={scores.engagement} />
        </div>
        <div className="mt-4 p-3 bg-gray-100 rounded-md text-center">
            <p className="text-sm text-gray-600 font-mono">
                Final Score = (Base × TypeCoefficient) + Freshness + Engagement
            </p>
            <p className="text-lg font-bold text-gray-800">
                {scores.final.toFixed(2)} = ({scores.base} × {scores.type_coefficient}) + {scores.freshness} + {scores.engagement}
            </p>
        </div>
    </div>
);

const RawPayloadViewer: React.FC<{ payload: Record<string, any> }> = ({ payload }) => {
    const [isCollapsed, setIsCollapsed] = useState(true);

    return (
        <div>
            <button
                onClick={() => setIsCollapsed(!isCollapsed)}
                className="text-lg font-semibold text-gray-700 w-full text-left flex items-center justify-between"
            >
                Raw Payload
                <svg xmlns="http://www.w3.org/2000/svg" className={`h-5 w-5 transition-transform ${isCollapsed ? 'rotate-0' : 'rotate-180'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
            </button>
            {!isCollapsed && (
                <pre className="mt-2 bg-gray-800 text-white p-4 rounded-md text-xs overflow-x-auto max-h-96">
                    <code>{JSON.stringify(payload, null, 2)}</code>
                </pre>
            )}
        </div>
    );
};


export const ContentDetailDrawer: React.FC<ContentDetailDrawerProps> = ({ contentId, onClose }) => {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['contentDetail', contentId],
    queryFn: () => getContentDetail(contentId!),
    enabled: !!contentId,
  });

  const content = data?.data;

  return (
    <Drawer isOpen={!!contentId} onClose={onClose} title={isLoading ? "Loading..." : content?.title ?? "Content Detail"}>
        {isLoading && <p>Loading details...</p>}
        {isError && <p className="text-red-500">Error: {(error as Error).message}</p>}
        {content && (
            <div className="space-y-6">
                <section>
                    <div className="flex justify-between items-start">
                        <div>
                             <h3 className="text-2xl font-bold text-brand-dark">{content.title}</h3>
                             <div className="text-sm text-gray-500 mt-1 space-x-2">
                                <span>Provider: <span className="font-semibold">{content.provider}</span></span>
                                <span>|</span>
                                <span>Published: <span className="font-semibold">{formatDateTime(content.published_at)}</span></span>
                             </div>
                        </div>
                        <span className="capitalize text-sm font-medium bg-brand-blue text-white px-3 py-1 rounded-full">{content.type}</span>
                    </div>
                    <div className="mt-4">
                        {content.tags.map(tag => <Badge key={tag}>{tag}</Badge>)}
                    </div>
                </section>
                
                <hr/>

                <section>
                    <h4 className="text-lg font-semibold text-gray-700 mb-2">Metrics</h4>
                    <MetricsDisplay metrics={content.metrics} type={content.type}/>
                </section>

                <hr/>

                <section>
                    <h4 className="text-lg font-semibold text-gray-700 mb-2">Score Breakdown</h4>
                    <ScoresDisplay scores={content.scores}/>
                </section>
                
                {content.raw_payload && (
                    <>
                        <hr />
                        <section>
                            <RawPayloadViewer payload={content.raw_payload} />
                        </section>
                    </>
                )}
            </div>
        )}
    </Drawer>
  );
};
