import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/core/ui/label_item_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/date_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/priority_widget.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/user_notifier_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import '../../../home/widgets/content/avatar_widget.dart';
import 'button_widget.dart';

class CartDetailsWidget extends ConsumerWidget {
  final CardItem card;
  const CartDetailsWidget({super.key, required this.card});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final users = ref.watch(userNotifierProvider).users;
    final currentUser = users.firstWhere((user) => user.id == card.userId);
    return Center(
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: Container(
            width: 400,
            height: 300,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
            decoration: BoxDecoration(
              color: AppColors.white26,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              children: [
                Row(
                  children: [
                    Text(card.title),
                    const Spacer(),
                    DateWidget(date: card.createdAt, color: AppColors.twitter),
                  ],
                ),
                const SizedBox(height: 10),
                Text(card.description),
                const SizedBox(height: 10),
                if (card.category != null)
                  Column(
                    children: [
                      LabelItemWidget(labelItem: card.category!),
                      const SizedBox(height: 10),
                    ],
                  ),
                if (card.priority != null)
                  Column(
                    children: [
                      PriorityWidget(priority: card.priority!),
                      const SizedBox(height: 10),
                    ],
                  ),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    AvatarWidget(imageUrl: currentUser.image),
                    const SizedBox(width: 5),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            '${currentUser.name} ${currentUser.lastName ?? ''}',
                            style: AppTextStyles.style13,
                            overflow: TextOverflow.ellipsis,
                            maxLines: 1,
                            softWrap: false,
                          ),
                          if (currentUser.role != null)
                            LabelItemWidget(labelItem: currentUser.role!),
                        ],
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    ButtonWidget(
                      title: 'Удалить',
                      color: AppColors.red,
                      onTap: () => _handleDelete(context, ref, card.id),
                    ),
                    const SizedBox(width: 10),
                    ButtonWidget(
                      title: 'Сохранить',
                      color: AppColors.blue,
                      onTap: () {},
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _handleDelete(BuildContext context, WidgetRef ref, String id) {
    ref.read(cardNotifierProvider.notifier).removeCard(id);
    Navigator.of(context).pop();
  }
}
