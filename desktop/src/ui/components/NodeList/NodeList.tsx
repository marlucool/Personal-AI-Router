// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

import 'overlayscrollbars/styles/overlayscrollbars.css'
import { memo, useMemo, useState } from 'react'
import {
    Dropdown,
    Flex,
    Stack,
    Text,
    TextInput,
    type DropdownEntry
} from '@nvidia/foundations-react-core'
import { OverlayScrollbarsComponent } from 'overlayscrollbars-react'
import type { PartialOptions } from 'overlayscrollbars'
import type { NodeItem } from '@/shared/types/nodes'
import { useOverviewNodes } from '@/ui/hooks/useOverviewNodes'
import NodeCardDetails from './NodeCardDetails'
import { CONNECTIONS_WIDTH } from '@/ui/constants/app'
import OfflineNode from './OfflineNode'
import { FilterList } from '@/ui/components/icons'

const SCROLLBAR_OPTIONS = {
    scrollbars: { autoHide: 'leave', autoHideDelay: 800 }
} satisfies PartialOptions

function matchesNode(node: NodeItem, normalizedQuery: string): boolean {
    if (!normalizedQuery) return true

    const searchableText = [
        node.name,
        node.ipAddress,
        node.os,
        node.topology.cpu.model,
        ...node.topology.gpus.map(gpu => gpu.name)
    ]
        .join(' ')
        .toLowerCase()

    return searchableText.includes(normalizedQuery)
}

function NodeList() {
    const allNodes = useOverviewNodes()
    const [query, setQuery] = useState('')
    const [showOffline, setShowOffline] = useState(true)
    const [gpuOnly, setGpuOnly] = useState(false)

    const normalizedQuery = query.trim().toLowerCase()

    const filteredNodes = useMemo(
        () =>
            allNodes.filter(node => {
                if (showOffline === false && node.status === 'offline') return false
                if (gpuOnly && node.topology.gpus.length === 0) return false
                return matchesNode(node, normalizedQuery)
            }),
        [allNodes, gpuOnly, normalizedQuery, showOffline]
    )

    const { online, offline } = useMemo(() => {
        const on: NodeItem[] = []
        const off: NodeItem[] = []

        filteredNodes.forEach(node => {
            if (node.status !== 'offline') {
                on.push(node)
            } else {
                off.push(node)
            }
        })

        return { online: on, offline: off }
    }, [filteredNodes])

    const filterItems: DropdownEntry[] = useMemo(
        () => [
            {
                kind: 'checkbox',
                checked: gpuOnly,
                onCheckedChange: checked => setGpuOnly(checked === true),
                children: (
                    <Text kind="body/regular/sm" className="whitespace-nowrap">
                        GPU nodes only
                    </Text>
                )
            },
            {
                kind: 'checkbox',
                checked: showOffline,
                onCheckedChange: checked => setShowOffline(checked === true),
                children: (
                    <Text kind="body/regular/sm" className="whitespace-nowrap">
                        Show offline nodes
                    </Text>
                )
            }
        ],
        [gpuOnly, showOffline]
    )

    if (allNodes.length === 0) {
        return <Stack className="grow min-w-0 h-full" />
    }

    const visibleNodeCount = online.length + offline.length

    return (
        <Stack className="grow min-w-0 h-full max-w-300">
            <Flex
                align="center"
                gap="2"
                className="min-w-0 w-full mb-2"
                data-node-list-toolbar
            >
                <TextInput
                    placeholder="Search nodes by name, IP, GPU or CPU"
                    value={query}
                    onValueChange={setQuery}
                    className="flex-1 min-w-0 max-h-[32px]"
                    aria-label="Search nodes"
                />
                <Dropdown
                    items={filterItems}
                    attributes={{
                        DropdownContent: {
                            className: 'no-drag-elements node-list-filter-dropdown-content'
                        }
                    }}
                >
                    <FilterList style={{ fontSize: 16 }} />
                    <Text kind="body/bold/sm">Filter</Text>
                </Dropdown>
            </Flex>

            {(normalizedQuery || gpuOnly || !showOffline) && (
                <Text
                    kind="body/regular/sm"
                    className="text-subtle-color mb-2 ml-1"
                    aria-live="polite"
                >
                    Showing {visibleNodeCount} of {allNodes.length} node
                    {allNodes.length === 1 ? '' : 's'}
                </Text>
            )}

            {visibleNodeCount === 0 ? (
                <Stack className="min-h-24 items-center justify-center">
                    <Text kind="body/regular/sm" className="text-subtle-color">
                        No nodes match the current search and filters.
                    </Text>
                </Stack>
            ) : (
                <OverlayScrollbarsComponent
                    className="node-list-scroll-container"
                    style={{
                        padding: `${CONNECTIONS_WIDTH / 2}px`,
                        margin: `0 -${CONNECTIONS_WIDTH / 2}px`
                    }}
                    options={SCROLLBAR_OPTIONS}
                    defer
                >
                    <Stack className="min-w-0 min-h-full dir-ltr" gap="3" data-node-list-content>
                        {online.map(node => (
                            <NodeCardDetails key={node.id} node={node} />
                        ))}

                        {offline.length > 0 &&
                            offline.map(node => (
                                <OfflineNode
                                    key={node.id}
                                    nodeId={node.id}
                                    name={node.name}
                                    ipAddress={node.ipAddress}
                                />
                            ))}
                    </Stack>
                </OverlayScrollbarsComponent>
            )}
        </Stack>
    )
}

export default memo(NodeList)
