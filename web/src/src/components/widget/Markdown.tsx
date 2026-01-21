import 'katex/dist/katex.min.css';
import '@fancyapps/ui/dist/fancybox/fancybox.css';

import { Fancybox } from '@fancyapps/ui';
import { useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeKatex from 'rehype-katex';
import rehypeRaw from 'rehype-raw';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';

import { CodeBlock } from '../widget/CodeBlock';

interface IMarkdown {
    readonly className?: string;
    readonly children: string;
}

export const Markdown = ({ className, children }: IMarkdown) => {
    useEffect(() => {
        Fancybox.bind('[data-viewer]');
        return () => {
            Fancybox.destroy();
        };
    }, []);

    return (
        <div className={className}>
            <ReactMarkdown
                className="prose-sm max-w-[100%] break-words"
                children={children}
                components={{
                    a: ({ ...props }) => (
                        <a
                            className="text-purple-600 hover:underline"
                            href={props.href}
                            target="_blank"
                            {...props}
                        >
                            {props.children}
                        </a>
                    ),
                    img: ({ ...props }) => (
                        <a
                            className="flex flex-col items-center"
                            href={props.src}
                            data-viewer
                            data-caption={props.alt}
                        >
                            <img className="rounded-lg shadow-lg" {...props} alt={props.alt} />
                            <span className="-mt-6 text-sm">{props.alt}</span>
                        </a>
                    ),
                    pre: ({ ...props }) => <pre className="bg-transparent p-2" {...props} />,
                    code: ({ className, children }) => {
                        const match = /language-(\w+)/.exec(className ?? '');
                        const lang = match !== null ? match[1] : '';
                        const code = String(children);
                        return match ? (
                            <CodeBlock language={lang} fileName={`snippet_${Date.now()}.txt`}>
                                {code}
                            </CodeBlock>
                        ) : (
                            <code className="overflow-scroll rounded-sm bg-gray-200 px-2 py-0.5 font-mono text-gray-700">
                                {code.replace(/\n$/, '')}
                            </code>
                        );
                    },
                    table: ({ ...props }) => (
                        <table className="block overflow-x-auto whitespace-nowrap" {...props} />
                    )
                }}
                urlTransform={(url) => url}
                rehypePlugins={[rehypeKatex, rehypeRaw]}
                remarkPlugins={[remarkGfm, remarkMath]}
            />
        </div>
    );
};
